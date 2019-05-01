<template>
   <div class="accession content">
      <h2>
         Archives Records Transfer System Admin Panel<span class="login"><b>Logged in as:</b>{{loginName}}</span>
      </h2>
      <div>
          <router-link to="/admin"><i class="fas fa-arrow-left"></i>&nbsp;Back to Accessions</router-link>
      </div>
   </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  name: 'accession',
  computed: {
      ...mapGetters({
         loginName: 'admin/loginName',
      })
   },
   created() {
      this.$store.dispatch("admin/getAccessionDetail", this.$route.params.id)
   }
}
</script>

<style scoped>
a {
   color: cornflowerblue;
   font-weight: 500;
   text-decoration: none;
}
a:hover {
   text-decoration: underline;
}
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