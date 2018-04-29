export default {
  data() {
    return {
      items: [],
    }
  },
  created: function () {
    this.fetchAccounts();
  },
  methods: {
    arrayLink(account_id) {
      return {to: {name: 'ArrayList',params: {account_id: account_id}}}
    },
    fetchAccounts() {
      {
          this.$http.get(`${process.env.API_URL}/api/accounts`).then((response) => {
            console.log(response.body);
          this.items = response.body;
          return
        });
      }
    },
  }
}
