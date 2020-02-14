Vue.component('stock-info', {
  props: ['team'],
  data: function() {
    return {
      loading: false,
      teamData: {},
    }
  },
  watch: {
    team: function(newTeam, oldTeam) {
      this.loading = true;
      ax.get('/teams/' + this.team).then(response => {
        this.teamData = response.data;
        this.loading = false;
      }).catch(error => {
        console.log(error);
        console.log(error.data);
      });
    },
  },
  template: `
    <div class="col8">
      <div v-if="loading">
        <div class="alert alert-info" role="alert">
          Fetching team...
        </div>
      </div>
      <div v-else>
        <div v-if="team == -1">
          Select a team from the stock ticker to view data.
        </div>
        <div v-else>
          Current Team ID: {{ team }} <br />
          Current Team: {{ teamData.Name }}
        </div>
      </div>
    </div>
  `
});
