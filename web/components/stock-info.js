Vue.component('stock-info', {
  props: ['team'],
  template: `
    <div class="col8">
      Current Team: {{ team }}
    </div>
  `
});
