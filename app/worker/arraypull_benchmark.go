package main

import (
	"log"
	"os"
	"sync"
	"time"
	"github.com/pkg/profile"
	"github.com/sjeanpierre/SJP_Go_Packages/lib/rightscale"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func s(rs rightscale.Client, x *sync.WaitGroup) {
	defer timeTrack(time.Now(), "non-parallel")
	ArrayList, err := rs.Arrays(false)
	if err != nil {
		log.Fatalln(err)
	}
	_ = ArrayList
	x.Done()
}

func p(rs rightscale.Client) {
	defer timeTrack(time.Now(), "Parallel")
	ArrayList, err := rs.ArraysParallel(false)
	if err != nil {
		log.Fatalln(err)
	}
	_ = ArrayList
}

func main() {
	defer profile.Start().Stop()
	rs, err := rightscale.New(os.Getenv("RS_REFRESH_TOKEN"), "https://us-3.rightscale.com")
	log.Println(rs.BearerToken)
	if err != nil {
		log.Fatal(err)
	}
	//wg.Add(1)
	//go s(rs,&wg)
	p(rs)

}
