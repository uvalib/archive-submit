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
        <div class="list-controls">
          <div class="search pure-button-group" role="group">
            <input type="text" id="search"><button class="search pure-button pure-button-primary">Search</button>
          </div>
          <AccessionPager/>
        </div>
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
               <td class="center"><span v-html="typeIcon(acc.physical)"></span></td>
               <td class="center"><span v-html="typeIcon(acc.digital)"></span></td>
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
import AccessionPager from "@/components/AccessionPager";
export default {
   name: "admin",
   components: {
     AccessionPager,
   },
   computed: {
      ...mapState({
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
      },
      typeIcon( hasType ) {
         if (hasType) {
            return '<i style="color:green;font-size:20px;opacity:0.6;" class="fas fa-check-circle"></i>'
         } else {
            return '<i style="color:firebrick;font-size:20px;opacity:0.6;" class="fas fa-times-circle"></i>'
         }
      }
   },
   created() {
      this.$store.dispatch("admin/getAccessions");
   }
};
</script>

<style scoped>
div.list-controls {
  position:relative;
  margin: 25px 0 5px 0;
}
div.search {
  font-size: 14px;
}
div.search button.search.pure-button {
  padding: 3px 15px;
}
div.search input {
  width:250px;
  outline:none;
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
div.error {
   font-style: italic;
   color: firebrick;
   padding: 5px 00 15px;
}
table {
   width: 100%;
   font-size: 0.85em;
   color: #444;
}
table td {
  border-bottom: 1px solid #ccc;
}
table td.center {
  text-align: center;
}
tr.accession:hover {
   background: #f5f5f5;
   cursor: pointer;
}
</style>