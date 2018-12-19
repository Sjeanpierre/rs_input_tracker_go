export default {
  created: function () {
    this.setVars();
    this.fetchInputData();
    this.fetchArrayDetails();
    this.fetchAccountDetails();
  },
  methods: {
    setVars() {
      let account = this.$route.params.account_id;
      let array = this.$route.params.array_id;
      let input = this.$route.params.input;
      if (typeof account !== 'undefined' && account) {
        this.account_id = account
      }
      if (typeof array !== 'undefined' && array) {
        this.array_id = array
      }
      if (typeof input !== 'undefined' && input) {
        this.input_name = input
      }

    },
    fetchArrayDetails() {
      {
        let account = this.$route.params.account_id
        let array = this.$route.params.array_id
        if (typeof account != 'undefined' && account && typeof array != 'undefined' && array) {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}/arrays/${array}`).then((response) => {
            console.log(response.body);
          this.array_name = response.body.array_name;
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
        if (typeof account != 'undefined' && account && typeof array != 'undefined' && array)
        {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}/arrays/${array}/inputs/${this.input_name}`).then((response) => {
            console.log(response.body);
          this.items = response.body;
          this.input_versions_loading = false;
          return
        });

        }
      }
    },
    fetchAccountDetails() {
      {
        let account = this.$route.params.account_id
        if (typeof account != 'undefined' && account) {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}`).then((response) => {
            console.log(response.body);
            this.account_name = response.body.name;
            return
          });
        }
      }
    },
    format_date_as_ago(date) {
      let rDate = date;
      if (typeof rDate !== 'undefined' && rDate) {
        var m = this.moment(rDate,"YYYY-MM-DD HH:mm:ssZ");
        return m.fromNow();
      }
    },
    format_date_as_ago_long(date) {
      let rDate = date;
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
        {icon: 'shuffle', text: 'Switch Accounts',link: { name: 'Home'}},
        {icon: 'keyboard_return', text: 'Back to Array',link: { name: 'ArrayInputs', params: { account_id: this.$route.params.account_id, array_id: this.$route.params.array_id }}},
        {divider: true},
      ],
      states: [],
      a1: null,
      search: '',
      array_name: "",
      account_name: "",
      account_id: "",
      array_count: "",
      input_name: "",
      input_versions_loading: true,
      headers: [
        {text: 'Version', value: 'version'},
        {text: 'Type', value: 'type'},
        {text: 'Retrieved', value: 'created_at'},
        {text: 'Value', value: 'value'}
      ],
      items: [],
    }
  },
}
