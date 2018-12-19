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
    lastUpdatedInDays(array) {
      let updated = Date.parse(array.updated_at)
      let now = Date.now()
      Math.abs(now - updated)
      let TimeDiff = Math.abs(now - updated)
      return Math.ceil(TimeDiff / (1000 * 3600 * 24))
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
          let elements = [];
          let ref = this; //I keep the this context prior to entering the loop, which would change the meaning of this
          //to reduce noise I want to only include arrays that have been audited in the past 5 days
          response.body.forEach(function (element) {
            if (ref.lastUpdatedInDays(element) < 5) {
              elements.push(element)
            }
          });
          if (elements.length !== 0) {
            this.items = elements;
          } else {
            this.items = response.body;
          }
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


