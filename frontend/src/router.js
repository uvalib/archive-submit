import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import Submit from './views/Submit.vue'
import Admin from './views/Admin.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    },
    {
      path: '/submit',
      name: 'Submit',
      component: Submit
    },{
      path: '/admin',
      name: 'Admin',
      component: Admin
    }
  ]
})
