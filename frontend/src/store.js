import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import { getField, updateField } from 'vuex-map-fields';

Vue.use(Vuex)

// root state object. Holds all of the state for the system
const state = {
  showInventory: false,
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
  accession: {
    identifier: null,
    summary: '',
    activities: '',
    creator: '',
    selectedGenres: [],
    accessionType: 'new'
  },
  digital: {
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
    hasSoftware: '0',
    inventory: []
  }
}

// state getter functions. All are functions that take state as the first param 
// and the getters themselves as the second param. Getter params are passed 
// as a function. Access as a property like: this.$store.getters.NAME
const getters = {
  getField,
  submissionID: state => {
    return state.accession.identifier
  },
  inventoryCount: state => {
    return state.physical.inventory.length
  },
  inventoryItem: (state) => (idx) => {
    return state.physical.inventory[idx]
  },
  hasError: state => {
    return state.error != null && state.error != ""
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
  clearInventory(state) {
    state.physical.inventory = []
  },
  addInventory(state) {
    state.physical.inventory.push({boxNum: '', recordGroup: '', title: '',
      description:'', dates: '' })
  },
  updateInventory(state, payload) {
    state.physical.inventory[payload.idx] = payload.item
  },
  deleteInventory(state, idx) {
    state.physical.inventory.splice(idx, 1)
  },
  toggleInventory(state) {
    state.showInventory = !state.showInventory
  },
  clearSubmissionData(state) {
    state.accession = { identifier: '', summary: '', activities: '', creator: '', selectedGenres: [],accessionType: 'new' }
    state.digital = { description: '', dateRange: '', selectedTypes: [], 
      uploadedFiles: [], totalSizeBytes: 0 }
    state.physical = { dateRange: '', boxInfo: '', selectedTypes: [], transferMethod: 0, hasDigital: '1',
      techInfo: '', mediaCarriers: [], mediaCount: '', hasSoftware: '0', inventory: [] }
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
  setSubmissionID (state, identifier) {
    state.accession.identifier = identifier
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
    ctx.commit('setGenres', []) 
    axios.get("/api/genres").then((response)  =>  {
      ctx.commit('setGenres', response.data )
    }).catch(() => {
      ctx.commit('setError', "Internal Error: Unable to get genres")
    })
  },
  getMediaCarriers( ctx ) {
    ctx.commit('setMediaCarrierChoices', []) 
    axios.get("/api/media-carriers").then((response)  =>  {
      ctx.commit('setMediaCarrierChoices', response.data )
    }).catch(() => {
      ctx.commit('setError', "Internal Error: Unable to get media carriers")
    })
  },
  getTransferMethods( ctx ) {
    ctx.commit('setTransferMethods', []) 
    axios.get("/api/transfer-methods").then((response)  =>  {
      ctx.commit('setTransferMethods', response.data )
    }).catch(() => {
      ctx.commit('setError', "Internal Error: Unable to get transfer methods")
    })
  },
  getRecordTypes( ctx ) {
    ctx.commit('setRecordTypes', []) 
    axios.get("/api/types").then((response)  =>  {
      ctx.commit('setRecordTypes', response.data )
    }).catch(() => {
      ctx.commit('setError', "Internal Error: Unable to get record types") 
    })
  },
  getSubmissionID( ctx ) {
    ctx.commit('setSubmissionID', "") 
    axios.get("/api/identifier").then((response)  =>  {
      ctx.commit('setSubmissionID', response.data )
    }).catch(() => {
      ctx.commit('setError', "Internal Error: Unable to get SubmissionID") 
    })
  },
  removeUploadedFile( ctx, file ) { 
    ctx.commit("removeUploadedFile",file.name)
    axios.delete("/api/upload/"+file.name+"?key="+ctx.getters.submissionID)
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