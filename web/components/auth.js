Vue.component('auth', {
    data: function() {
        return {
            user: '',
        }
    },
    template: `
        <div v-if user != ''>
            <strong>User</strong>
            <p>Sign In</p>
            <p>Sign Up</p>
        </div><div v-else>
            <strong>{{ user }}</strong>
            <p>Sign Out</p>
        </div>
    `
})