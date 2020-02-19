Vue.component('team-split-info', {
  props: ['team', 'split'],
  template: `
    <div>
      <h5>{{ split.League }} {{ split.Season }} {{ split.Year }}</h5>
      <p>show stock data for {{ team.Name }}</p>
    </div>
  `
});
