
const Dashboard = { 
    template: '#dashboard' ,
    created: function (){
        this.$http.get('/api/routes').then(res => {
            this.routes = res.data
        })
        this.$http.get('/api/config').then(res => {
            this.config = res.data
        }).catch(err => {
            console.log(err)
        })
    },
    data : function () {
        return {
            title: 'Dashboard',
            routes: {},
            config: {},
            search:""
        }
    },
    methods: {},
    computed: {
        routesFiltered (){
            var self = this
            var col = _.map(this.routes, function (v,k,o) {
                v.ID = k
                return v
            })
            return _.filter(col, function (v, k, o) {
                return v.ID.toLowerCase().includes(self.search.toLowerCase())
            })
        }
    }
}
const Routes = { 
    template: '#routes' ,
    created: function () {
        this.$http.get('/api/routes').then(res => {
            this.routes = res.data
        })
    },
    data : function () {
        return {
            title: 'Routes',
            routes: {}
        }
    },
    methods: {
        toggleRoute: function (id) {
            this.routes[id].Enabled = !this.routes[id].Enabled
            this.$http.put(`/api/routes/${id}`, this.routes[id]).then(res => {
                this.$http.get(`/api/routes`).then(resp => {
                    this.routes = resp.data
                })
            })
        },
    },
    computed: {}
}
const Route = { 
    template: '#route' ,
    created: function () {
        this.$http.get(`/api/routes/${this.$route.params.id}`).then(res => {
            this.route = res.data
            this.app = this.$route.params.id
        })
    },
    data : function () {
        return {
            title: 'Route',
            route: {
                Stream: "",
                CopyKey: false,
                Enabled: true,
                Endpoints: []
            },
            app: "",
            newEndpoint: {
                Host: '',
                Port: "1935",
                App: '',
                Stream: '',
                Enabled: true,
            },
            editingEndpoint: {
                Host: '',
                Port: "1935",
                App: '',
                Stream: '',
                Enabled: true,
            },
            editingIndex: 0
        }
    },
    methods: {
        showAddEndpoint: function() {
            $('#addEndpoint').modal('show')
        },
        toggleEndpoint: function (i) {
            this.route.Endpoints[i].Enabled = !this.route.Endpoints[i].Enabled
        },
        addEndpoint: function (){
            if(this.route.Endpoints == null || typeof this.route.Endpoints == "undefined"){
                this.route.Endpoints = []
            }
            this.route.Endpoints.push(Vue.util.extend({}, this.newEndpoint))
            $('#addEndpoint').modal('hide')
        },
        editEndpoint: function (i) {
            this.editingIndex = i
            this.editingEndpoint = Vue.util.extend({}, this.route.Endpoints[i])
            $('#editEndpoint').modal('show')
        },
        saveEditEndpoint: function (){
            $('#editEndpoint').modal('hide')
            this.route.Endpoints[this.editingIndex] = this.editingEndpoint
            
        },
        removeEndpoint: function(i){
            this.route.Endpoints.splice(i,1)
        },
        save: function () {
            this.$http.put(`/api/routes/${this.$route.params.id}`, this.route).then(res => {
                this.$router.push('/routes')
            })
        }
    },
    computed: {}
}
const NewRoute = { 
    template: '#route' ,
    created: function () {
    },
    data : function () {
        return {
            title: 'Route',
            route: {
                Stream: '',
                CopyKey: false,
                Enabled: true,
                Endpoints: []
            },
            app: "",
            newEndpoint: {
                Host: '',
                Port: "1935",
                App: '',
                Stream: '',
                Enabled: true,
            },
            editingEndpoint: {
                Host: '',
                Port: "1935",
                App: '',
                Stream: '',
                Enabled: true,
            },
            editingIndex: 0
        }
    },
    methods: {
        showAddEndpoint: function() {
            $('#addEndpoint').modal('show')
        },
        toggleEndpoint: function (i) {
            this.route.Endpoints[i].Enabled = !this.route.Endpoints[i].Enabled
        },
        addEndpoint: function (){
            this.route.Endpoints.push(Vue.util.extend({}, this.newEndpoint))
            $('#addEndpoint').modal('hide')
        },
        editEndpoint: function (i) {
            this.editingIndex = i
            this.editingEndpoint = Vue.util.extend({}, this.route.Endpoints[i])
            $('#editEndpoint').modal('show')
        },
        saveEditEndpoint: function () {
            $('#editEndpoint').modal('hide')
            this.route.Endpoints[this.editingIndex] = this.editingEndpoint

        },
        removeEndpoint: function(i){
            this.route.Endpoints.splice(i,1)
        },
        save: function () {
            this.$http.post(`/api/routes`, this.route).then(res => {
                this.$router.push('/routes')
            })
        }
    },
    computed: {}
}
const Configuration = { 
    template: '#configuration' ,
    created: function (){
        this.$http.get('/api/config').then(res => {
            this.config = res.data
        }).catch(err => {
            console.log(err)
        })
    },
    data : function () {
        return {
            title: 'Configuration',
            config:{}
        }
    },
    methods: {},
    computed: {}
}
const Stats = { 
    template: '#stats' ,
    created: {},
    data : function () {
        return {
            title: 'Stats'
        }
    },
    methods: {},
    computed: {}
}
const Admin = { 
    template: '#admin' ,
    created: {},
    data : function () {
        return {
            title: 'Admin'
        }
    },
    methods: {},
    computed: {}
}
Vue.component('btnToolbar', {
    template: '#btn-toolbar',
    props:["title"]
})
Vue.component('bCard', {
    template: "#card",
    props:['title', 'content', 'header']
})
const routes = [
    { path: '/', redirect: '/dashboard'},
    { path: '/dashboard', component: Dashboard },
    { path: '/routes', component: Routes },
    { path: '/routes/new', component: NewRoute },
    { path: '/routes/:id', component: Route },
    { path: '/configuration', component: Configuration },
    { path: '/stats', component: Stats },
    { path: '/admin', component: Admin },
]
const router = new VueRouter({
    routes, // short for `routes: routes`
    linkActiveClass: "active"
})
Vue.prototype.$http = axios
const app = new Vue({
    router: router,
}).$mount("#app")