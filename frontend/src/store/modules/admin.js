import axios from 'axios'

const admin = {
   namespaced: true,
   state: {
      accessions: [],
      totalAccessions: 0,
      filteredTotal: 0,
      page: 0,
      pageSize: 0,
      accessionDetail: null,
      notes: [],
      queryStr: "",
      tgtGenre: "",
      addingNote: false,
      working: false
   },
   getters: {
      loginName(_state, _getters, rootState) {
         if (rootState.user == null) {
            return ""
         }
         return rootState.user.firstName + " (" + rootState.user.email + ")"
      },
      hasNotes(state) {
         return state.notes.length > 0
      }
   },
   mutations: {
      setWorking(state, val) {
         state.working = val
      },
      setAddingNote(state, adding) {
         state.addingNote = adding
      },
      setGenreFilter(state, val) {
         state.tgtGenre = val
      },
      resetAccessionsSearch(state) {
         state.tgtGenre = ""
         state.queryStr = ""
         state.page = 1
         state.filteredTotal = 0
         state.total = 0
         state.accessions.length = 0
      },
      setAccessionsPage(state, resp) {
         state.totalAccessions = resp.total
         state.filteredTotal = resp.filteredTotal
         state.page = resp.page
         state.pageSize = resp.pageSize
         state.accessions = resp.accessions
      },
      clearAccessionDetail(state) {
         state.accessionDetail = null
      },
      setAccessionDetail(state, data) {
         state.accessionDetail = data
      },
      setNotes(state, data) {
         state.notes = data
      },
      addNote(state, note) {
         state.notes.push(note)
      },
      updateSearchQuery(state, val) {
         state.queryStr = val
      },
      gotoFirstPage(state) {
         state.page = 1
      },
      gotoLastPage(state) {
         state.page = Math.floor(state.totalAccessions / state.pageSize) + 1
      },
      nextPage(state) {
         state.page++
      },
      prevPage(state) {
         state.page--
      },
   },
   actions: {
      firstPage(ctx) {
         ctx.commit('gotoFirstPage')
         ctx.dispatch("getAccessionsPage")
      },
      prevPage(ctx) {
         ctx.commit('prevPage')
         ctx.dispatch("getAccessionsPage")
      },
      nextPage(ctx) {
         ctx.commit('nextPage')
         ctx.dispatch("getAccessionsPage")
      },
      lastPage(ctx) {
         ctx.commit('gotoLastPage')
         ctx.dispatch("getAccessionsPage")
      },
      getAccessionsPage(ctx) {
         ctx.commit("setLoading", true, { root: true })
         let url = "/api/admin/accessions?page=" + ctx.state.page
         if (ctx.state.queryStr.length > 0) {
            url = url + "&q=" + ctx.state.queryStr
         }
         if (ctx.state.tgtGenre.length > 0 ) {
            url = url +"&g="+ctx.state.tgtGenre
          }
         axios.get(url, { withCredentials: true }).then((response) => {
            ctx.commit('setAccessionsPage', response.data)
            ctx.commit("setLoading", false, { root: true })
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get accessions", { root: true })
            ctx.commit("setLoading", false, { root: true })
         })
      },
      getAccessionDetail(ctx, id) {
         ctx.commit("setLoading", true, { root: true })
         ctx.commit('clearAccessionDetail')
         axios.get("/api/admin/accessions/" + id, { withCredentials: true }).then((response) => {
            ctx.commit('setAccessionDetail', response.data)
            ctx.commit("setLoading", false, { root: true })
            ctx.dispatch('getAccessionNotes', id)
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get accession detail", { root: true })
            ctx.commit("setLoading", false, { root: true })
         })
      },
      getAccessionNotes(ctx, id) {
         axios.get("/api/admin/accessions/" + id+"/notes", { withCredentials: true }).then((response) => {
            ctx.commit('setNotes', response.data)
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get accession notes", { root: true })
         })
      },
      addNote(ctx, data) {
         let id = ctx.state.accessionDetail.id
         data.userID = ctx.rootState.user.id
         ctx.commit("setWorking", true)
         axios.post("/api/admin/accessions/" + id+"/notes", data, { withCredentials: true }).then((response) => {
            ctx.commit('addNote', response.data)
            ctx.commit("setWorking", false)
            ctx.commit('setAddingNote', false)
         }).catch((err) => {
            ctx.commit('setError', err.response.data, { root: true })
            ctx.commit("setWorking", false)
         })
      }
   }
}
export default admin