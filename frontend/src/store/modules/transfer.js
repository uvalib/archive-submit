
import axios from 'axios'
import { getField, updateField } from 'vuex-map-fields'

const transfer = {
   namespaced: true,

   state: {
      showInventory: false,
      genres: [],
      digitalTransfer: true,
      physicalTransfer: false,
      digitalRecordTypes: [],
      physicalRecordTypes: [],
      mediaCarrierChoices: [],
      transferMethods: [],
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
         transferMethod: "1",
         hasDigital: '1',
         techInfo: '',
         mediaCarriers: [],
         mediaCount: '',
         hasSoftware: '0',
         inventory: []
      }
   },

   getters: {
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
      digitalUploadSize: state => {
         let mb = state.digital.totalSizeBytes/(1000.0*1000.0)
         if (mb > 1000.0) {
            return mb/1000.0.toFixed(2)+"GB"
         } else {
            return mb.toFixed(2)+"MB"
         }
      }, 
   },

   mutations: {
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
      clearRecordTypes(state) {
         state.digitalRecordTypes = []
         state.physicalRecordTypes = []
      }
   },

   actions: {
      getGenres( ctx ) {
         ctx.commit('setGenres', []) 
         axios.get("/api/genres").then((response)  =>  {
            ctx.commit('setGenres', response.data )
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get genres", {root: true})
         })
      },
      getMediaCarriers( ctx ) {
         ctx.commit('setMediaCarrierChoices', []) 
         axios.get("/api/media-carriers").then((response)  =>  {
            ctx.commit('setMediaCarrierChoices', response.data )
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get media carriers", {root: true})
         })
      },
      getTransferMethods( ctx ) {
         ctx.commit('setTransferMethods', []) 
         axios.get("/api/transfer-methods").then((response)  =>  {
            ctx.commit('setTransferMethods', response.data )
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get transfer methods", {root: true})
         })
      },
      getRecordTypes( ctx ) {
         ctx.commit('clearRecordTypes', []) 
         axios.get("/api/types").then((response)  =>  {
            ctx.commit('setRecordTypes', response.data )
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get record types", {root: true})
         })
      },
      getSubmissionID( ctx ) {
         ctx.commit('setSubmissionID', "") 
         axios.get("/api/identifier").then((response)  =>  {
            ctx.commit('setSubmissionID', response.data )
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get SubmissionID", {root: true}) 
         })
      },
      removeUploadedFile( ctx, file ) { 
         ctx.commit("removeUploadedFile",file.name)
         axios.delete("/api/upload/"+file.name+"?key="+ctx.getters.submissionID)
      }
   }
}
export default transfer