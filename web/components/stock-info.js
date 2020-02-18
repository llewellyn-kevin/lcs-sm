Vue.component('stock-info', {
  props: ['team'],
  data: function() {
    return {
      loading: false,
      teamData: {},
      splits: {},
    }
  },
  watch: {
    team: function(newTeam, oldTeam) {
      this.loading = true;
      ax.get('/teams/' + this.team).then(response => {
        this.teamData = response.data;
      }).catch(error => {
        console.log(error);
        console.log(error.data);
      }).finally(() => {
        this.teamData.StockValues.forEach(stock => {
          if(this.splits["split-" + stock.SplitID] == null) {
            this.splits["split-" + stock.SplitID] = 'loading';
            ax.get('/splits/' + stock.SplitID).then(response => {
              console.log(response.data);
              this.splits['split-' + stock.SplitID] = response.data;
            }).catch(error => {
              console.log(error);
              console.log(error.data);
            });
          }
          console.log(this.splits);
        });
        this.loading = false;
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
          <h3>{{ teamData.Name }}</h3>
          <p v-for="split in splits">{{ split.League }} {{ split.Season}} {{ split.Year }}</p>
        </div>
      </div>
    </div>
  `
});
