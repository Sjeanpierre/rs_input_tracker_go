export default {
  data() {
    return {
      search: '',
      array_name: "",
      account_name: "",
      array_count: "",
      array_list_loading: true,
      headers: [
        {text: 'Name', value: 'array_name'},
        {text: 'Array ID', value: 'array_id',align: 'center'},
        {text: 'Rightscale Link',align: 'center', sortable: false}
      ],
      items: [],
    }
  },
  created: function () {
    this.fetchArrayList();
    this.fetchAccountDetails();
  },
  methods: {
    fetchArrayDetails() {
      {
        let account = this.$route.params.account_id;
        let array = this.$route.params.array_id;
        if (typeof account !== 'undefined' && account && typeof array !== 'undefined' && array) {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}/arrays/${array}`).then((response) => {
            console.log(response.body);
            this.array_name = response.body.array_name;
        });
        }
      }
    },
    fetchInputData() {
      {
        let account = this.$route.params.account_id;
        let array = this.$route.params.array_id;
        // I saw some instances where undefined route vars were coming through to the request, this guards against that
        if (typeof account !== 'undefined' && account && typeof array !== 'undefined' && array)
        {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}/arrays/${array}/inputs`).then((response) => {
            console.log(response.body);
            this.items = response.body;
        });

        }
      }
    },
    fetchArrayList() {
      {
        let account = this.$route.params.account_id;
        // I saw some instances where undefined route vars were coming through to the request, this guards against that
        if (typeof account !== 'undefined' && account)
        {
          this.$http.get(`${process.env.API_URL}/api/accounts/${account}/arrays`).then((response) => {
            console.log(response.body);
          this.items = response.body;
          this.array_list_loading = false;
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
    rightscaleLink(array_id) {
      let account = this.$route.params.account_id;
      let url = `https://my.rightscale.com/acct/${account}/server_arrays/${array_id}`
      return {
        href: url,
        target: '_blank',
        rel: 'noopener'
      }
    }
  }
}


