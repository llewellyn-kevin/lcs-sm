Vue.component('team-split-info', {
  props: ['team', 'split'],
  data: function() {
    return {
      stocks: [],
    };
  },
  created: function() {
    ax.get('/splits/' + this.split.ID + '/teams/' + this.team.ID + '/stock-values').then(response => {
      this.stocks = response.data 
    }).catch(error => {
      console.log(error);
      console.log(error.data);
    });
  },
  template: `
    <div>
      <h5>{{ split.League }} {{ split.Season }} {{ split.Year }}</h5>
      <ul class="list-group">
        <li 
          v-for="stock in stocks"
          class="list-group-item list-group-item-action d-flex justify-content-between align-items-center">
          <strong>Week {{ stock.Week }}</strong>
          <span>{{ stock.Value }}</span>
        </li>
      </ul>
    </div>
  `
});
