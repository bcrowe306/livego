
const Dashboard = { 
    template: '#dashboard' ,
    created: function (){
        this.$http.get('http://localhost:8080/api/routes').then(res => {
            this.routes = res.data
        })
    },
    data : function () {
        return {
            title: 'Dashboard',
            routes: {}
        }
    },
    methods: {},
    computed: {}
}
const Routes = { 
    template: '#routes' ,
    created: function () {
        this.$http.get('http://localhost:8080/api/routes').then(res => {
            this.routes = res.data
        })
    },
    data : function () {
        return {
            title: 'Routes',
            routes: {}
        }
    },
    methods: {},
    computed: {}
}
const Route = { 
    template: '#route' ,
    created: function () {
        this.$http.get(`http://localhost:8080/api/routes/${this.$route.params.id}`).then(res => {
            this.route = res.data
            console.log(res.data)
            this.app = this.$route.params.id
        })
    },
    data : function () {
        return {
            title: 'Routes',
            route: {},
            app: ""
        }
    },
    methods: {},
    computed: {}
}
const Configuration = { 
    template: '#configuration' ,
    created: {},
    data : function () {
        return {
            title: 'Configuration'
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
    template: `
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pb-2 mb-3 border-bottom">
        <h1 class="h2">{{title}}</h1>
        <div class="btn-toolbar mb-2 mb-md-0">
            <div class="btn-group mr-2">
                <button class="btn btn-sm btn-outline-secondary">Share</button>
                <button class="btn btn-sm btn-outline-secondary">Export</button>
            </div>
            <button class="btn btn-sm btn-outline-secondary dropdown-toggle">
                <span data-feather="calendar"></span>
                This week
            </button>
        </div>
    </div>`,
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