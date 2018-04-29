export default {
  created: function () {
    this.setVars();
    this.fetchInputData();
    this.fetchArrayDetails();
    this.fetchArrayCount();
    this.fetchAccountDetails();
  },
  methods: {
    setVars() {
      let account = this.$route.params.account_id;
      let array = this.$route.params.array_id;
      if (typeof account !== 'undefined' && account) {
        this.account_id = account
      }
      if (typeof array !== 'undefined' && array) {
        this.array_id = array
      }

    },
    fetchArrayDetails() {
      {
        let account = this.$route.params.account_id;
        let array = this.$route.params.array_id;
        if (typeof account !== 'undefined' && account && typeof array !== 'undefined' && array) {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}/arrays/${array}`).then((response) => {
            console.log(response.body);
          this.array_name = response.body.array_name;
          this.array_details = response.body;
          return
        });
        }
      }
    },
    fetchInputData() {
      {
        let account = this.$route.params.account_id
        let array = this.$route.params.array_id
        // I saw some instances where undefined route vars were coming through to the request, this guards against that
        if (typeof account !== 'undefined' && account && typeof array !== 'undefined' && array)
        {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}/arrays/${array}/inputs`).then((response) => {
            console.log(response.body);
          this.items = response.body;
          this.datatable_loading = false;
        });

        }
      }
    },
    fetchArrayCount() {
      {
        let account = this.$route.params.account_id;
        // I saw some instances where undefined route vars were coming through to the request, this guards against that
        if (typeof account !== 'undefined' && account)
        {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}/arrays`).then((response) => {
            console.log(response.body);
          this.array_count = response.body.length
        });

        }
      }
    },
    fetchAccountDetails() {
      {
        let account = this.$route.params.account_id;
        if (typeof account !== 'undefined' && account) {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}`).then((response) => {
            console.log(response.body);
          this.account_name = response.body.name;
        });
        }
      }
    },
    inputVersionLink(input_name) {
      return {to: {name: 'InputVersions',params: {account_id: this.account_id,array_id: this.array_id,input: input_name}}}
    },
  },
  computed: {
    array_last_updated_date() {
      let rDate = this.array_details.updated_at;
      if (typeof rDate !== 'undefined' && rDate) {
        var m = this.moment(rDate,"YYYY-MM-DD HH:mm:ssZ");
        return m.fromNow();
      }
    },
    array_last_updated_date_long() {
      let rDate = this.array_details.updated_at;
      if (typeof rDate !== 'undefined' && rDate) {
        var m = this.moment(rDate,"YYYY-MM-DD HH:mm:ssZ");
        return m.calendar();
      }
    },
  },
  data() {
    return {
      drawer: null,
      menu_items: [
        {divider: true},
        {heading: 'Actions'},
        {snackbarToogle: true},
        {icon: 'shuffle', text: 'Switch Accounts',link: { name: 'Home'}},
        {icon: 'keyboard_return', text: 'Switch Array',link: { name: 'ArrayList', params: { account_id: this.$route.params.account_id }}},
        {divider: true},
      ],
      states: [],
      snackbar: false,
      a1: null,
      search: '',
      array_name: "",
      account_name: "",
      account_id: "",
      array_count: "",
      array_details: "",
      datatable_loading: true,
      headers: [
        {text: 'Name', value: 'name'},
        {text: 'Type', value: 'type'},
        {text: 'Version', value: 'version'},
        {text: 'Value', value: 'value'}
      ],
      items: [],
    }
  },
}
