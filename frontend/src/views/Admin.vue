<template>
   <div class="admin index">
      <h2>
         System Admin Panel
         <span class="login">
            <b>Logged in as:</b>
            {{loginName}}
         </span>
      </h2>
      <div>
        <div class="list-controls">
         <div class="search pure-button-group" role="group">
            <input @input="updateSearchQuery" @keyup.enter="searchClicked" type="text" id="search" :value="queryStr">
            <button  @click="searchClicked"  class="search pure-button pure-button-primary">Search</button>
         </div>
         <span class="tag-filter" v-if="tgtGenre.length > 0">
            <b>Genre:</b> {{tgtGenre}} <i @click="removeFilter" class="unfilter fas fa-times-circle"></i>
         </span>  
         <AccessionPager/>
        </div>
         <table class="pure-table">
            <thead>
               <th>Identifier</th>
               <th>Type</th>
               <th>Submitter</th>
               <th>Description</th>
               <th>Genres</th>
               <th style="width:50px">Physical</th>
               <th style="width:50px">Digital</th>
               <th style="width:50px">Notes</th>
               <th>Transferred</th>
            </thead>
            <tr v-for="acc in accessions" :key="acc.id" class="accession" :data-id="acc.id" @click="accessionClicked">
               <td>{{ acc.accessionID }}</td>
               <td>{{ acc.type }}</td>
               <td>{{ acc.submitter }}</td>
               <td>{{ acc.description }}</td>
               <td>
                  <span class="tag" v-for="(tag,idx) in tagList(acc)" :key="idx" @click="tagClicked">{{tag}}</span> 
               </td>
               <td class="center"><span v-html="typeIcon(acc.physical)"></span></td>
               <td class="center"><span v-html="typeIcon(acc.digital)"></span></td>
               <td class="center">{{ acc.notesCount }}</td>
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
         error: state => state.error,
         loading: state => state.loading,
         queryStr: state => state.admin.queryStr,
         tgtGenre: state => state.admin.tgtGenre,
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
      },
      updateSearchQuery(e) {
         this.$store.commit('admin/updateSearchQuery', e.target.value)
      },
      searchClicked() {
         this.$store.dispatch("admin/getAccessionsPage")
      },
      tagList( acc ) {
         if (acc.genres) {
            return acc.genres.split(",")
         } 
         return []
      },
      tagClicked(event) {
         event.stopPropagation()
         // textContent may return whitepace before/after tag. Strip it
         let tag = event.currentTarget.textContent.replace(/^\s+|\s+$/g, '')
         this.$store.commit('admin/setGenreFilter', tag)
         this.$store.dispatch("admin/getAccessionsPage")
      },
      removeFilter() {
         this.$store.commit('admin/setGenreFilter', "")
         this.$store.dispatch("admin/getAccessionsPage")
      }
   },
   created() {
      this.$store.commit("admin/resetAccessionsSearch");
      this.$store.dispatch("admin/getAccessionsPage");
   }
};
</script>

<style scoped>
div.admin {
   background: white;
   color: #444;
   position: relative;
   min-width: 1000px;
   padding: 30px 50px 250px 50px;
}
.tag-filter {
   font-size: 0.8em;
   border: 1px solid #ccc;
   padding: 2px 4px 0 10px;
   border-radius: 20px;
   cursor: pointer;
}
span.tag {
   display: inline-block;
   margin: 0 4px 4px 0;
   font-size: 0.9em;
   background: #0078e7;
   color: white;
   padding: 2px 10px 2px 10px;
   font-weight: 500;
   opacity: 0.6;
   border-radius: 20px;
}
span.tag:hover {
   cursor: pointer;
   opacity: 1;
}
i.fas.unfilter {
   color: firebrick;
   margin-left: 5px;
   opacity: 0.6;
}
i.fas.unfilter:hover {
   cursor: pointer;
   opacity: 1;
}
div.list-controls {
  position:relative;
  margin: 25px 0 5px 0;
}
div.search {
  font-size: 14px;
  display: inline-block;
  margin-right: 10px;
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