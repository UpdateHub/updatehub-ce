import Vue from 'vue'
import App from './App.vue'
import Axios from 'axios'
import VueRouter from 'vue-router'

import Login from './components/Login'
import Overview from './components/Overview'
import DeviceList from './components/DeviceList'
import DeviceDetails from './components/DeviceDetails'
import PackageList from './components/PackageList'
import PackageDetails from './components/PackageDetails'
import RolloutList from './components/RolloutList'
import RolloutNew from './components/RolloutNew'
import RolloutDetails from './components/RolloutDetails'

function errorResponseHandler(error) {
  if (error.response.status === 401) {
    app.currentUser = null;
    router.push("/login?redirect=" + "/")
  }
}

Axios.interceptors.response.use(
  response => response,
  errorResponseHandler
);

Vue.config.productionTip = false;

Vue.prototype.$http = Axios

Vue.use(VueRouter)

const routes = [
  { path: '/', redirect: '/overview' },
  { path: '/login', component: Login },
  { path: '/overview', component: Overview },
  { path: '/devices', component: DeviceList },
  { path: '/devices/:uid', component: DeviceDetails },
  { path: '/packages', component: PackageList },
  { path: '/packages/:uid', component: PackageDetails },
  { path: '/rollouts', component: RolloutList },
  { path: '/rollouts/new', component: RolloutNew },
  { path: '/rollouts/:id', component: RolloutDetails }
]

const router = new VueRouter({ routes })

var app = new Vue({
  render: h => h(App),

  router,

  computed: {
    currentUser: {
      cache: false,

      get () {
        if (!localStorage.currentUser) {
          return null
        } else {
          return JSON.parse(localStorage.currentUser)
        }
      },

      set (currentUser) {
        localStorage.currentUser = JSON.stringify(currentUser)
        Vue.prototype.$http.defaults.headers.common['Authorization'] =
          'Bearer ' + currentUser.token
      }
    }
  }
})

if (app.currentUser) {
  Vue.prototype.$http.defaults.headers.common['Authorization'] =
    'Bearer ' + app.currentUser.token
}

Vue.prototype.$app = app

app.$mount('#app')
