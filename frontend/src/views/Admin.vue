<template>
   <div class="admin content">
      <h2>
         Archives Records Transfer System Admin Panel
         <span class="login">
            <b>Logged in as:</b>
            {{loginName}}
         </span>
      </h2>
      <div>
         <table class="pure-table">
            <thead>
               <th>Identifier</th>
               <th>Type</th>
               <th>Submitter</th>
               <th>Description</th>
               <th>Genres</th>
               <th>Physical</th>
               <th>Digital</th>
               <th>Transferred</th>
            </thead>
            <tr v-for="acc in accessions" :key="acc.id" class="accession" :data-id="acc.id" @click="accessionClicked">
               <td>{{ acc.accessionID }}</td>
               <td>{{ acc.type }}</td>
               <td>{{ acc.submitter }}</td>
               <td>{{ acc.description }}</td>
               <td>{{ acc.genres }}</td>
               <td>{{ acc.physical }}</td>
               <td>{{ acc.digital }}</td>
               <td>{{ acc.submittedAt.split("T")[0] }}</td>
            </tr>
         </table>
      </div>

      <div class="error">{{error}}</div>
   </div>
</template>

<script>
import { mapState } from "vuex";
import { mapGetters } from "vuex";
export default {
   name: "admin",
   computed: {
      ...mapState({
         total: state => state.admin.totalAccessions,
         page: state => state.admin.page,
         accessions: state => state.admin.accessions,
         error: state => state.error
      }),
      ...mapGetters({
         loginName: "admin/loginName"
      })
   },
   methods: {
      accessionClicked(event) {
         let tgt = event.currentTarget;
         let accID = tgt.dataset.id;
         this.$router.push("admin/accessions/" + accID);
      }
   },
   created() {
      this.$store.dispatch("admin/getAccessions");
   }
};
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
table {
   width: 100%;
   font-size: 0.85em;
   color: #444;
   margin-top:25px;
}
tr.accession:hover {
   background: #f5f5f5;
   cursor: pointer;
}
</style>