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

      var ctx = this.$refs['graph-canvas'].getContext('2d');
      var weekNumbers = [];
      var values = [];

      this.stocks.forEach(stock => {
        weekNumbers.push(stock.Week);
        values.push(stock.Value);
      });

      new Chart(ctx, {
        "type": "line",
        "data": {
          "labels": weekNumbers,
          "datasets": [{
            "label": this.team.Name + " Stock Trends",
            "data": values,
            "fill": false,
            "borderColor": "rgb(75, 192, 192)",
            "lineTension": 0.1 
          }]
        },
        "options": {
          "legend": {
            "display": false
          }
        }
      });

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

      <canvas ref="graph-canvas" width="800" height="600"></canvas>
    </div>
  `
});
