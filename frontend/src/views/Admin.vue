<template>
  <div class="admin content">
    <h2>Archives Records Transfer System Admin Panel<span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
export default {
   name: "admin",
   computed: {
      ...mapState({
         total: state => state.admin.totalAccessions,
         page: state => state.admin.page,
         accessions: state => state.admin.submissions,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
      })
   },
   methods: {
  
   },
   created() {
      let authUser = this.$cookies.get("archives_xfer_user")
      if (authUser) {
         this.$store.commit("setUser", authUser)
         this.$store.dispatch("admin/getSubmissions")
      } else {
         this.$router.push("forbidden")
      }
   }
}
</script>

<style scoped>
span.login {
   font-family: sans-serif;
   font-size: 0.6em;
   float: right;
   font-weight: 100;
   color: #666;
}
span.login b {
   margin-right: 5px;
}
</style>