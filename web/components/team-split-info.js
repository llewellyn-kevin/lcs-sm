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
        var mostRecent = this.stocks.reduce((a, i) => { return i.Week > a.Week ? {"Week":i.Week,"Value":i.Value} : a }, zero);
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
            },
            "responsive": true,
            "maintainAspectRatio": true,
            "aspectRatio": 1.77
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
      <div class="row">

        <div class="col-md">
          <ul class="list-group weekly-stocks-list">
            <li 
              v-for="stock in stocks"
              class="list-group-item list-group-item-action d-flex justify-content-between align-items-center">
              <strong>Week {{ stock.Week }}</strong>
              <span>{{ stock.Value }}</span>
            </li>
          </ul>
        </div>

        <div class="col-md">
          <canvas ref="graph-canvas"></canvas>

          <form>
            <div v-if="stockCreated" class="alert alert-success">
              Stock Succesfully Created
            </div>
            
            <div class="row">
              <div class="col-sm"><label>Week</label></div>
              <div class="col-sm"><input v-model="newWeek"></div>
              <div class="w-100"></div>
              <div class="col-sm"><label>Value</label></div>
              <div class="col-sm"><input v-model="newValue"></div>
              <div class="w-100"></div><div class="col-sm"></div>
              <div class="col-sm">
                <button type="button" class="btn btn-primary" v-on:click="createStock">Create Stock</button>
              </div></div>
          </form>
        </div>

      </div>
    </div>
  `
});
