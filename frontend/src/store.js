import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import { getField, updateField } from 'vuex-map-fields';

Vue.use(Vuex)

// root state object. Holds all of the state for the system
const state = {
  genres: [],
  error: null,
  isUVA: false,
  digitalTransfer: true,
  physicalTransfer: false,
  digitalRecordTypes: [],
  physicalRecordTypes: [],
  mediaCarrierChoices: [],
  transferMethods: [],
  user: {
    id: 0,
    firstName: '',
    lastName: '',
    email: '',
    phone: '',
    title: '',
    affiliation: ''
  },
  general: {
    summary: '',
    activities: '',
    creator: '',
    selectedGenres: [],
    accessionType: 'new'
  },
  digital: {
    uploadID: null,
    description: '',
    dateRange: '',
    selectedTypes: [],
    uploadedFiles: [],
    totalSizeBytes: 0
  },
  physical: {
    dateRange: '',
    boxInfo: '',
    selectedTypes: [],
    transferMethod: 0,
    hasDigital: '1',
    techInfo: '',
    mediaCarriers: [],
    mediaCount: '',
    hasSoftware: '0'
  }
}

// state getter functions. All are functions that take state as the first param 
// and the getters themselves as the second param. Getter params are passed 
// as a function. Access as a property like: this.$store.getters.NAME
const getters = {
  getField,
  uploadID: state => {
    return state.digital.uploadID
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
  isUVA: state => {
    return state.isUVA
  },
  user: state => {
    return state.user
  },
  digitalTransfer: state => {
    return state.digitalTransfer
  },
  physicalTransfer: state => {
    return state.physicalTransfer
  },
  physicalRecordTypes: state => {
    return state.physicalRecordTypes
  },
  digitalRecordTypes: state => {
    return state.digitalRecordTypes
  },
  mediaCarrierChoices: state => {
    return state.mediaCarrierChoices
  },
  transferMethods: state => {
    return state.transferMethods
  },
  digitalUploadSize: state => {
    let mb = state.digital.totalSizeBytes/(1000.0*1000.0)
    if (mb > 1000.0) {
      return mb/1000.0.toFixed(2)+"GB"
    } else {
      return mb.toFixed(2)+"MB"
    }
  }, 
}

// Synchronous updates to the state. Can be called directly in components like this:
// this.$store.commit('mutation_name') or called from asynchronous actions
const mutations = {
  updateField,
  clearPhysicalXefrDigitalInfo(state) {
    state.physical.techInfo =  ''
    state.physical.mediaCarriers = []
    state.physical.mediaCount = ''
    state.physical.hasSoftware = '0'
  },
  setGenres (state, genres) {
    state.genres = genres
  },
  setMediaCarrierChoices (state, carriers) {
    state.mediaCarrierChoices = carriers
  },
  setTransferMethods (state, methods) {
    state.transferMethods = methods
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
    state.user = {firstName: "", lastName:"", title:"", affiliation:"", email:email, phone:""}
  },
  clearUser(state) {
    state.user = {id: 0, firstName: "", lastName:"", title:"", affiliation:"", email:"", phone:""}
  },
  setError (state, error) {
    state.error = error
  },
  setUploadID (state, uploadID) {
    state.digital.uploadID = uploadID
  },
  addUploadedFile (state, file) {
    state.digital.uploadedFiles.push(file.name)
    state.digital.totalSizeBytes += file.size
  },
  removeUploadedFile (state, file) {
    let index = state.digital.uploadedFiles.indexOf(file.name)
    if (index !== -1) {
      state.digital.uploadedFiles.splice(index, 1)
      state.digital.totalSizeBytes -= file.size
    }
  },
  setDigitalTransfer(state, digital) {
    state.digitalTransfer = digital
  },
  setPhysicalTransfer(state, physical) {
    state.physicalTransfer = physical
  },
  setRecordTypes(state, types) {
    types.forEach(function( rt ) {
      if (rt.digitalOnly) {
        state.digitalRecordTypes.push(rt)
      } else {
        state.physicalRecordTypes.push(rt)
      }
    });
  },
}

// Actions are asynchronous calls that commit mutatations to the state.
// All actions get context as a param which is essentially the entirety of the 
// Vuex instance. It has access to all getters, setters and commit. They are 
// called from components like: this.$store.dispatch('action_name', data_object)
const actions = {
  getGenres( ctx ) {
    axios.get("/api/genres").then((response)  =>  {
      ctx.commit('setGenres', response.data )
    }).catch(() => {
      ctx.commit('setGenres', []) 
      ctx.commit('setError', "Internal Error: Unable to get genres")
    })
  },
  getMediaCarriers( ctx ) {
    axios.get("/api/media-carriers").then((response)  =>  {
      ctx.commit('setMediaCarrierChoices', response.data )
    }).catch(() => {
      ctx.commit('setMediaCarrierChoices', []) 
      ctx.commit('setError', "Internal Error: Unable to get media carriers")
    })
  },
  getTransferMethods( ctx ) {
    axios.get("/api/transfer-methods").then((response)  =>  {
      ctx.commit('setTransferMethods', response.data )
    }).catch(() => {
      ctx.commit('setTransferMethods', []) 
      ctx.commit('setError', "Internal Error: Unable to get transfer methods")
    })
  },
  getRecordTypes( ctx ) {
    axios.get("/api/types").then((response)  =>  {
      ctx.commit('setRecordTypes', response.data )
    }).catch(() => {
      ctx.commit('setRecordTypes', []) 
      ctx.commit('setError', "Internal Error: Unable to get record types") 
    })
  },
  getUploadID( ctx ) {
    axios.get("/api/identifier").then((response)  =>  {
      ctx.commit('setUploadID', response.data )
    }).catch(() => {
      ctx.commit('setUploadID', "") 
      ctx.commit('setError', "Internal Error: Unable to get uploadID") 
    })
  },
  removeUploadedFile( ctx, file ) { 
    ctx.commit("removeUploadedFile",file.name)
    axios.delete("/api/upload/"+file.name+"?key="+ctx.getters.uploadID)
  }
}

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
  state,
  getters,
  actions,
  mutations,
  plugins: [errorPlugin] 
})