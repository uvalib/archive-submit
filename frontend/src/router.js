import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import Access from './views/Access.vue'
import Submit from './views/Submit.vue'
import Thanks from './views/Thanks.vue'
import Admin from './views/Admin.vue'
import Accession from './views/Accession.vue'
import Forbidden from './views/Forbidden.vue'
import Verify from './views/Verify.vue'
import store from './store'

Vue.use(Router)

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/access',
      name: 'access',
      component: Access
    },
    {
      path: '/submit',
      name: 'submit',
      component: Submit,
      beforeEnter: (to, from, next) => {
        // must have a user in the store or in the cookie to access the submit form
        // Note: can also call Vue.cookies.keys() to get an array of all cookie names
        // this could be used to look fo a _shibsession* cookie in the app domain to
        // ensure authorization
        let authUser = Vue.cookies.get("archives_xfer_user")
        if (store.state.user == null && authUser == null) {
          next("/")
        } else {
          next()
        }
      }
    },
    {
      path: '/thanks',
      name: 'thanks',
      component: Thanks
    },
    {
      path: '/verify/:token',
      name: 'verify',
      component: Verify
    },
    {
      path: '/admin',
      name: 'admin',
      component: Admin,
      meta: { requiresAuth: true }
    },
    {
      path: '/admin/accessions/:id',
      name: 'accession',
      component: Accession,
      meta: { requiresAuth: true }
    },
    {
      path: '/forbidden',
      name: 'forbidden',
      component: Forbidden
    }
  ],
  scrollBehavior(/*to, from, savedPosition*/) {
    return { x: 0, y: 0 }
  },
})

router.beforeEach((to, _from, next) => {
  if (to.meta.requiresAuth == true) {
    let getters = store.getters
    let isAuth = getters["isAuthenticated"]
    if ( isAuth === false) {
      let authUser = Vue.cookies.get("archives_xfer_user")
      if (authUser) {
        authUser.authenticated = true
        store.commit("setUser", authUser)
        Vue.cookies.remove("archives_xfer_user");
        Vue.cookies.remove("archives_xfer_settings");
      } else {
        window.location.href = "/authenticate?url=" + to.fullPath
      }
    }
  } 
  next()
})

export default router
