var teamDirectory = new Vue({
    el: '#team-directory',
    data: {
        teams: [],
    },
    created: () => {
        ax.get('/teams').then(response => {
            this.teams = response;
        }).catch(error => {
            console.log(error);
            console.log(error.response);
        })
    },
});