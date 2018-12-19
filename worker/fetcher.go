package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sjeanpierre/SJP_Go_Packages/lib/rightscale"
	"github.com/sjeanpierre/rs_input_tracker_go/app/models"
)

type inputAudit struct {
	Array       rightscale.ServerArray
	ArrayInputs rightscale.Inputs
	Account     string
}

func perform(rsToken, rsAccountID, rsEndpoint string) {
	rs, err := rightscale.New(rsToken, rsEndpoint)
	//log.Println(rs.BearerToken)
	//os.Exit(0)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Getting array list")
	ArrayList, err := rs.ArraysParallel(false)
	if err != nil {
		log.Fatalln(err)
	}
	//open new DB connection here to avoid thundering herd from the go routines
	//if we wait until we are down in the go routine to open the new connection then they will all do so simultaneously
	log.Println("Establishing DB connection")
	_ = models.ListAccounts()
	acct := rsAccountID
	var wg sync.WaitGroup
	for _, a := range ArrayList {
		wg.Add(1)
		go func(account string, array rightscale.ServerArray, x *sync.WaitGroup) {
			ia := inputAudit{Array: array, Account: account}
			populateArrayInputs(rs, &ia)
			log.Printf("Array %s Contains %v inputs", array.Name, len(ia.ArrayInputs))
			writeInputsToDB(&ia)
			x.Done()
		}(acct, a, &wg)
	}
	wg.Wait()
	a := stringToINT(acct)
	populateArrays(ArrayList, a)
}

func populateArrayInputs(rs rightscale.Client, ia *inputAudit) {
	array := ia.Array
	inputs, err := rs.ArrayInputs(array)
	if err != nil {
		log.Fatal(err)
	}
	ia.ArrayInputs = inputs
}

func writeInputsToDB(ia *inputAudit) {
	arrayID, err := ia.Array.ArrayID()
	if err != nil {
		log.Printf("Could not process Server array %s, invalid Href '%s'", ia.Array.Name, ia.Array.Href)
		return
	}
	//get list of latest inputs to allow value compare - using account_id,array_id
	//also includes inactive inputs to prevent creating new series for ones that were previously present then inactivated
	currentInputList := models.ListCurrentInputsWithInactive(ia.Account, arrayID)
	inputMap := make(map[string]models.Input)
	for _, input := range currentInputList {
		inputMap[input.Name] = input
	}
	//todo, if we ever find that we are spending a lot of time in SQL with the inserts, we can
	//change them to be performed in bulk fashion which should reduce the overhead a bit
	//In order to do this we'd collect all createInput function calls into a single Inputs struct
	//then pass that in the the DB logic as a group
	for _, newInput := range ia.ArrayInputs {
		//check for existence of new input
		input, ok := inputMap[newInput.Name]
		if !ok {
			//input by this name does not yet exist, we can create it and exit this iteration of the loop
			array := stringToINT(arrayID)
			accID := stringToINT(ia.Account)
			ni := models.Input{ArrayID: int(array), AccountID: int(accID), RawValue: newInput.Value, Name: newInput.Name, Version: 1}
			models.CreateInput(ni, true)
			continue
		}
		if input.Inactive {
			//Previous series of input had been marked as inactive
			//input is now present reactivating the series
			models.ReactivateInput(input)
		}
		if input.RawValue != newInput.Value {
			//input exists,but the value has changed from the last audited value, increment version and insert new record
			array := stringToINT(arrayID)
			accID := stringToINT(ia.Account)
			v := input.Version + 1
			newInput := models.Input{ArrayID: int(array), AccountID: int(accID), RawValue: newInput.Value, Name: newInput.Name, Version: v}
			models.CreateInput(newInput, true)
			continue
		} else {
			//input was found and value was not changed.
			//log.Printf("Skipping: %s Reason: unchanged", newInput.Name)
		}
	}
	//determine which records have been removed and should be marked as dead records
	//add tombstone record with incremented version
	//if the latest version of an input is a tombstone record, then omit it from the results
	missingInputs := findMissingInputs(ia.ArrayInputs,currentInputList)
	for _, missingInput := range missingInputs {
		log.Printf("Notice: Input %s is missing from array %s",missingInput.Name,ia.Array.Name)
		v := missingInput.Version + 1
		newInput := models.Input{ArrayID: missingInput.ArrayID, AccountID: missingInput.AccountID,
			RawValue: missingInput.RawValue, Name: missingInput.Name, Version: v}
		models.CreateInactiveInputRecord(newInput, true)
	}
}

//determine which records have been removed and should be marked as dead records
//compare the current inputs against the new set pulled in
//check if an input exists in the current set, but does not exist in the new set
//create new record for that input with same values as last highest version
//set deleted flag to true
//all functions seeking to determine latest arrays should omit records with deleted = true
func findMissingInputs(new rightscale.Inputs, old models.Inputs) models.Inputs {
	var missing models.Inputs
	newInputMap := make(map[string]rightscale.Input) //stores newInput map for quick lookup
	for _, newInput := range new {
		newInputMap[newInput.Name] = newInput
	}
	for _, oldInput := range old {
		_, ok := newInputMap[oldInput.Name]
		if !ok && !oldInput.Inactive { //if missing and active
			missing = append(missing,oldInput)
		}
	}
	return missing
}

func populateArrays(arrays rightscale.ServerArrays, account int) {
	var exists rightscale.ServerArrays
	var nouveau rightscale.ServerArrays
	currentArrays := models.ListArraysByAccount(account)
	//split arrays into exist, and nouveau buckets existing
	for _, array := range arrays {
		arrayID, err := array.ArrayID()
		aID := stringToINT(arrayID)
		if err != nil {
			log.Println("Could not retrieve array ID for", array.Name)
		}
		if !currentArrays.Exists(account, int(aID)) {
			nouveau = append(nouveau, array)
			continue
		}
		exists = append(exists, array)
	}

	//action for nouveau Arrays, insert them

	for _, array := range nouveau {
		arr, _ := array.ArrayID()
		a := models.Array{AccountID: account,
			ArrayID: stringToINT(arr),
			Name:    array.Name}
		models.CreateArray(a)
	}
	//action for existing Arrays, update them
	//get db version of the exists set
	var a models.Arrays
	for _, array := range exists {
		arr, _ := array.ArrayID()
		n := currentArrays.FindArrayByAccountArrayID(account, stringToINT(arr))
		a = append(a, n)
	}
	models.UpdateArrays(a)
}

func stringToINT(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0
	}
	return int(i)
}

func Handler() (Response,error) {
	rsToken, tokenOK := os.LookupEnv("RS_REFRESH_TOKEN")
	rsAccount, accountOK := os.LookupEnv("RS_ACCOUNT_ID")
	rsEndpoint,endpointOK := os.LookupEnv("RS_ACCOUNT_ENDPOINT")
	if tokenOK && accountOK && endpointOK {
		perform(rsToken, rsAccount, rsEndpoint)
	} else {
		log.Fatalf("ENV var missing.\nPresent:\nRS_REFRESH_TOKEN" +
			" %v\nRS_ACCOUNT_ID %v\n RS_ACCOUNT_ENDPOINT %v ", tokenOK, accountOK,endpointOK)
	}

	return Response{Message:"Audit complete"},nil

}

type Response struct {
	Message string `json:"message"`
}

func main() {
	lambda.Start(Handler)
}
