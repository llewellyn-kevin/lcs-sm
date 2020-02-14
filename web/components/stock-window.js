Vue.component('stock-window', {
  data: function() {
    return {
      selectedTeam: -1,
    }
  },
  methods: {
    onChangeTeam: function(newTeam) {
      this.selectedTeam = newTeam;
    },
  },
  template: `
    <div class="row">
      <team-directory 
        v-on:change-team="onChangeTeam">
      </team-directory>
      <stock-info 
        v-bind:team="selectedTeam">
      </stock-info>
    </div>
  `
});
