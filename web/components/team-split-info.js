Vue.component('team-split-info', {
  props: ['team', 'split'],
  data: function() {
    return {
      stocks: [],
      newWeek: 0,
      newValue: 0,
      stockCreated: false
    };
  },
  methods: {
    createStock(e) {
      ax.post('/splits/'+this.split.ID+'/teams/'+this.team.ID+'/stock-values?week='+this.newWeek+'&value='+this.newValue)
      .then(response => {
        this.refreshStocks();

        this.stockCreated = true;
        setTimeout(function() {
          this.stockCreated = false;
        }, 1600);
      }).catch(error => {
        console.log(error);
        console.log(error.data);
      });
    },
    refreshStocks() {
      ax.get('/splits/' + this.split.ID + '/teams/' + this.team.ID + '/stock-values').then(response => {
        this.stocks = response.data;
        var zero = {"Week":0,"Value":0}
        var mostRecent = this.stocks.reduce((a, i) => {
          return i.Week > a.Week ? {"Week":i.Week,"Value":i.Value} : a 
        }, zero);
        this.newWeek = mostRecent.Week + 1;
        this.newValue = mostRecent.Value;

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
    }
  },
  created: function() {
    this.refreshStocks();
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

      <div v-if="stockCreated" class="alert alert-success">
        Stock Succesfully Created
      </div>
      <input v-model="newWeek">
      <input v-model="newValue">
      <button type="button" class="btn btn-primary" v-on:click="createStock">New Stock</button>
    </div>
  `
});
