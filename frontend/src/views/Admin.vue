<template>
  <div class="admin content">
    <h2>Archives Records Transfer System Admin Panel<span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
    <div>
         <h3>Accessions</h3>
         <table class="pure-table">
            <thead>
               <th style="width:12px"/>
               <th>Identifier</th>
               <th>Transferred</th>
               <th>Actions</th>
            </thead>
            <tr v-for="acc in accessions" :key="acc.id">
               <td><input :data-id="acc.id" type="checkbox"/></td>
               <td>{{ acc.accessionID }}</td>
               <td>{{ acc.submittedAt.split("T")[0] }}</td>
               <td>
                  <i title="view" class="action fas fa-eye"></i>
                  <i title="edit" class="action fas fa-edit"></i>
               </td>
            </tr>
         </table>
      </div>

    <div class="error">{{error}}</div>
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
         accessions: state => state.admin.accessions,
         error: state => state.error,
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
         this.$store.dispatch("admin/getAccessions")
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
div.error {
  font-style: italic;
  color: firebrick;
  padding: 5px 00 15px;
}
i.fas.action {
   cursor: pointer;
   opacity: 0.5;
   margin: 0 5px;
   font-size: 1.1em;
}
i.fas:hover {
   opacity: 1;
}
table {
   width: 100%;
   font-size: 0.85em;
   color: #444;
}
th.checkbox {
   width: 40px;
}
td.centered {
   text-align: center;
}
</style>