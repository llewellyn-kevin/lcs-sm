Vue.component('stock-info', {
  props: ['team'],
  data: function() {
    return {
      loading: false,
      teamData: {},
      splits: [],
      selectedSplit: 0,
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
        ax.get('/teams/' + this.team + '/splits/').then(response => {
          this.splits = response.data;
          this.splits.sort((a, b) => {
            // Sort by most recent Year first
            if(a.Year > b.Year) { return -1; }
            else if(a.Year < b.Year) { return 1; }
            else {
              // In the same year sort by Summer, then Spring
              if(a.Season == b.Season) { return 0; }
              else { return a.Season == "Summer" ? -1 : 1 }
            }
          });
          this.selectedSplit = this.splits[0]
          this.loading = false;
        }).catch(error => {
          console.log(error);
          console.log(error.data);
        });
      });
    },
  },
  template: `
    <div class="col-lg-9">
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
          <strong>Current Value: {{ teamData.CurrentValue }}</strong>
          <br />
          <select v-model="selectedSplit">
            <option v-for="split in splits" v-bind:value="split">
              {{ split.League }} {{ split.Season }} {{ split.Year }}
            </option>
          </select>
          <team-split-info 
            v-bind:key="selectedSplit.ID"
            v-bind:team="teamData" 
            v-bind:split="selectedSplit">
          </team-split-info>
        </div>
      </div>
    </div>
  `
});
