import axios from 'axios'

const admin = {
   namespaced: true,
   state: {
      accessions: [],
      totalAccessions: 0,
      page: 0,
      pageSize: 0,
      accessionDetail: null
   },
   getters: {
      isAuthenticated(state) {
         if (state.user == null) {
            return false
         }
         return state.user.authenticated
      },
      loginName (_state, _getters, rootState) {
         if (rootState.user == null) {
            return ""
         }
         return rootState.user.firstName + " ("+rootState.user.email+")"
      }
   },
   mutations: {
      setAccessions (state, accessionsInfo) {
         state.totalAccessions = accessionsInfo.total 
         state.page = accessionsInfo.page
         state.pageSize = accessionsInfo.pageSize
         state.accessions.length = 0
         accessionsInfo.accessions.forEach( function(acc) {
            state.accessions.push(acc)
         })
      },
      clearAccessions (state) {
         state.accessions.length = 0
         state.totalAccessions = 0
         state.page = 0
      },
      clearAccessionDetail (state) {
         state.accessionDetail = null
      },
      setAccessionDetail (state, data) {
         state.accessionDetail = data
      }
   },
   actions: {
      getAccessions( ctx ) {
         ctx.commit('clearAccessions')
         axios.get("/api/admin/accessions").then((response)  =>  {
            ctx.commit('setAccessions', response.data )
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get accessions", {root:true})
         })
      },
      getAccessionDetail(ctx, id) {
         ctx.commit('clearAccessionDetail')
         axios.get("/api/admin/accessions/"+id).then((response)  =>  {
            ctx.commit('setAccessionDetail', response.data )
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get accession detail", {root:true})
         })
      }
   }
}
export default admin