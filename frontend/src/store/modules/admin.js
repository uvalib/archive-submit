import axios from 'axios'

const admin = {
   state: {
      accessions: [],
      totalAccessions: 0,
      page: 0,
      pageSize: 0
   },
   mutations: {
      setAccessions (state, accessionsInfo) {
         state.totalAccessions = accessionsInfo.total 
         state.page = accessionsInfo.page
         state.pageSize = accessionsInfo.pageSize
         accessionsInfo.thumbs.forEach( function(acc) {
            state.accessions.push(acc)
         })
      },
      clearAccessions (state) {
         state.accessions = []
         state.totalAccessions = 0
         state.page = 0
      },
   },
   actions: {
      getGenres( ctx ) {
         ctx.commit('clearAccessions')
         axios.get("/api/admin/accessions").then((response)  =>  {
            ctx.commit('setAccessions', response.data )
         }).catch(() => {
            ctx.commit('setError', "Internal Error: Unable to get accessions", {root:true})
         })
      },
   }
}
export default admin