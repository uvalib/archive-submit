import Vue from 'vue'
import Vuex from 'vuex'
import admin from './modules/admin'
import transfer from './modules/transfer'
import { getField, updateField } from 'vuex-map-fields'

Vue.use(Vuex)

// Plugin to listen for error messages being set. After a delay, clear them
const errorPlugin = store => {
  store.subscribe((mutation) => {
    if (mutation.type === "setError") {
      if ( mutation.payload != null ) {
        setTimeout( ()=>{ store.commit('setError', null)}, 10000)
      }
    }
  })
}

// A Vuex instance is created by combining state, getters, actions and mutations
export default new Vuex.Store({
  // Global state/mutations/actions: for stuff that is used across modules
  state: {
    isUVA: false,
    error: null,
    user: null,
    loading: false
  },
  getters: {
    getField,
    hasError: state => {
      return state.error != null
    }
  },
  mutations: {
    updateField,
    setError (state, error) {
      state.error = error
    },
    setLoading (state, isLoading) {
      state.loading = isLoading
    },
    setUVA (state, isUVA) {
      state.isUVA = isUVA
    },
    setUser (state, user) {
      user.affiliation = user.affiliation.replace( /\+/g, ' ' )
      user.title = user.title.replace( /\+/g, ' ' )
      state.user = user
    },
    setUserEmail(state, email) {
      state.user = {email:email}
    },
    clearUser(state) {
      state.user = null
    },
  },
  modules: {
    admin: admin,
    transfer: transfer,
  },
  plugins: [errorPlugin] 
})