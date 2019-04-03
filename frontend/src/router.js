import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import Submit from './views/Submit.vue'
import Admin from './views/Admin.vue'
import Forbidden from './views/Forbidden.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/submit',
      name: 'submit',
      component: Submit
    },
    {
      path: '/admin',
      name: 'admin',
      component: Admin
    },
    {
      path: '/forbidden',
      name: 'forbidden',
      component: Forbidden
    }
  ]
})
