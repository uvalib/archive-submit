import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

// root state object. Holds all of the state for the system
const state = {
  genres: [],
  uploadID: null,
  error: null,
  uploadedFiles: [],
  isUVA: false,
  user: null
}

// state getter functions. All are functions that take state as the first param 
// and the getters themselves as the second param. Getter params are passed 
// as a function. Access as a property like: this.$store.getters.NAME
const getters = {
  uploadID: state => {
    return state.uploadID
  },
  genres: state => {
    return state.genres
  },
  error: state => {
    return state.error
  },
  hasError: state => {
    return state.error != null && state.error != ""
  },
  uploadedFiles: state => {
    return state.uploadedFiles
  },
  isUVA: state => {
    return state.isUVA
  },
  user: state => {
    return state.user
  },
}

// Synchronous updates to the state. Can be called directly in components like this:
// this.$store.commit('mutation_name') or called from asynchronous actions
const mutations = {
  setGenres (state, genres) {
    if (genres) {
      state.genres = genres
    }
  },
  setUVA (state, isUVA) {
    state.isUVA = isUVA
  },
  setUser (state, user) {
    user.affiliation = user.affiliation.replace( /\+/g, ' ' )
    user.title = user.title.replace( /\+/g, ' ' )
    state.user = user
  },
  setError (state, error) {
    state.error = error
  },
  setUploadID (state, uploadID) {
    state.uploadID = uploadID
  },
  addUploadedFile (state, filename) {
    state.uploadedFiles.push(filename)
  },
  removeUploadedFile (state, filename) {
    let index = state.uploadedFiles.indexOf(filename)
    if (index !== -1) {
      state.uploadedFiles.splice(index, 1)
    }
  }
}

// Actions are asynchronous calls that commit mutatations to the state.
// All actions get context as a param which is essentially the entirety of the 
// Vuex instance. It has access to all getters, setters and commit. They are 
// called from components like: this.$store.dispatch('action_name', data_object)
const actions = {
  getGenres( ctx ) {
    axios.get("/api/genres").then((response)  =>  {
      if ( response.status === 200) {
        ctx.commit('setGenres', response.data )
      } else {
        ctx.commit('setGenres', []) 
        ctx.commit('setError', "Internal Error: "+response.data) 
      }
    }).catch(() => {
      ctx.commit('setGenres', []) 
      ctx.commit('setError', "Internal Error: Unable to reach any services") 
    })
  },
  getUploadID( ctx ) {
    axios.get("/api/identifier").then((response)  =>  {
      if ( response.status === 200) {
        ctx.commit('setUploadID', response.data )
      } else {
        ctx.commit('setUploadID', []) 
        ctx.commit('setError', "Internal Error: "+response.data) 
      }
    }).catch(() => {
      ctx.commit('setUploadID', []) 
      ctx.commit('setError', "Internal Error: Unable to reach any services") 
    })
  },
  removeUploadedFile( ctx, filename ) {
    ctx.commit("removeUploadedFile",filename)
    axios.delete("/api/upload/"+filename+"?key="+ctx.getters.identifier)
  }
}

// Plugin to listen for error messages being set. After a delay, clear them
const errorPlugin = store => {
  store.subscribe((mutation) => {
    if (mutation.type === "setError") {
      if ( mutation.payload != null ) {
        setTimeout( ()=>{ store.commit('setError', null)}, 6000)
      }
    }
  })
}

// A Vuex instance is created by combining state, getters, actions and mutations
export default new Vuex.Store({
  state,
  getters,
  actions,
  mutations,
  plugins: [errorPlugin] 
})