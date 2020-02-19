Vue.component('team-directory', {
  data: function() {
    return {
      loading: true,
      teams: [],
    }
  },
  created: function() {
    ax.get('/teams').then(response => {
      this.teams = response.data;
      this.loading = false;
    }).catch(error => {
      console.log(error);
      console.log(error.data);
    });
  },
  template: `
    <div class="col4">
      <h4>Stock Ticker</h4>
      <div v-if="loading">
        <div class="alert alert-info" role="alert">
          Fetching teams...
        </div>
      </div>
      <div v-else>
        <div v-if="teams.length == 0">
          <div class="alert alert-danger" role="alert">
            No teams found.
          </div>
        </div>
        <div v-else>
          <div class="list-group">
            <button 
              v-for="team in teams" 
              v-on:click="$emit('change-team', team.ID)"
              class="list-group-item list-group-item-action d-flex justify-content-between align-items-center">
              <strong>{{ team.Name }}</strong>
              <span>{{ team.CurrentValue }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  `
});
