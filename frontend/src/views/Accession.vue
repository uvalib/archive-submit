<template>
   <div class="accession admin">
      <template v-if="loading">
         <h1>Loading...</h1>
      </template>
      <template v-else>
         <h2>
            System Admin Panel<span class="login"><b>Logged in as:</b>{{loginName}}</span>
         </h2>
         <router-link to="/admin"><i class="fas fa-arrow-left"></i>&nbsp;Back to Accessions</router-link>
         <div class="accession-details">
            <AccordionContent title="General Information" expanded>
               <div class="info-block">
                  <div><b>Transfer Identifier:</b><p>{{details.identifier}}</p></div>
                  <div><b>Transferred On</b><p>{{formattedDate(details.createdAt)}}</p></div>
                  <div><b>Accession Type:</b><p>{{details.accessionType}}</p></div>
                  <div><b>Summary:</b><p>{{details.summary}}</p></div>
                  <div><b>Activities Leading to Creation:</b><p>{{details.activities}}</p></div>
                  <div><b>Creator:</b><p>{{details.creator}}</p></div>
                  <div><b>Genres:</b><p>{{safeCSV(details.genres)}}</p></div>
               </div>
            </AccordionContent>
            <template v-if="details.digitalTransfer">
               <AccordionContent title="Digital Transfer Details">
                  <div class="info-block">
                     <div><b>Technical Description:</b><p>{{details.digital.description}}</p></div>
                     <div><b>Date Range of Files:</b><p>{{details.digital.dateRange}}</p></div>
                     <div><b>Record Types:</b><p>{{safeCSV(details.digital.selectedTypes)}}</p></div>
                     <div><b>Total Transfer Size:</b><p>{{(details.digital.totalSizeBytes/1000.0/1000.0).toFixed(2)}}GB</p></div>
                     <div><b>Files Transferred:</b><p>{{safeCSV(details.digital.uploadedFiles)}}</p></div>
                  </div>
               </AccordionContent>
            </template>
            <template v-if="details.physicalTransfer">
               <AccordionContent title="Physical Transfer Details">
                  <div class="info-block">
                     <div><b>Date Range of Records:</b><p>{{details.physical.dateRange}}</p></div>
                     <div><b>Number and Size of Boxes:</b><p>{{details.physical.boxInfo}}</p></div>
                     <div><b>Record Types:</b><p>{{safeCSV(details.physical.selectedTypes)}}</p></div>
                     <div><b>Transfer Method:</b><p>{{details.physical.transferMethodName}}</p></div>
                     <div><b>Transfer Method:</b><p>{{details.physical.transferMethodName}}</p></div>
                  </div>
               </AccordionContent>
               <template v-if="details.physical.hasDigital">
                  <AccordionContent title="Digital Media Carrier Details">
                     <div class="info-block">
                        <div><b>Technical Description:</b><p>{{details.physical.techInfo}}</p></div>
                        <div><b>Media Carriers:</b><p>{{safeCSV(details.physical.mediaCarriers)}}</p></div>
                        <div><b>Media Carrier Estimates:</b><p>{{details.physical.mediaCount}}</p></div>
                        <div><b>Transfer Includes Software:</b><p>{{details.physical.hasSoftware}}</p></div>
                     </div>
                  </AccordionContent>
               </template>
               <AccordionContent title="Inventory">
                  <div class="info-block">
                     <table class="pure-table" style="font-size:0.8em;width:100%;"> 
                        <thead>
                           <tr>
                              <th>Box Number</th>
                              <th>Records Group #</th>
                              <th style="width:175px">Box Title</th>
                              <th style="width:400px">Contents Description</th>
                              <th>Dates</th>
                           </tr>
                        </thead>
                        <tr v-for="(item, idx) in details.physical.inventory" :key="idx">
                           <td>{{item.boxNum}}</td>
                           <td>{{item.recordGroup}}</td>
                           <td>{{item.title}}</td>
                           <td>{{item.description}}</td>
                           <td>{{item.dates}}</td>
                        </tr>
                     </table>
                  </div>
               </AccordionContent>
            </template>
            <template v-if="hasNotes">
               <AccordionContent title="Notes">
                  <div class="note" v-for="(note, idx) in notes" :key="idx">
                     <p><b>{{note.title}}</b></p>
                     <span class="note-time"><b>{{note.userName}}</b>{{formattedDate(note.createdAt)}}</span>
                     <div>{{note.note}}</div>
                  </div>
               </AccordionContent>
            </template>
         </div>
      </template>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
import AccordionContent from '@/components/AccordionContent'
export default {
  name: 'accession',
  components: {
      AccordionContent: AccordionContent,
  },
  computed: {
      ...mapState({
         total: state => state.admin.totalAccessions,
         details: state => state.admin.accessionDetail,
         notes: state=>state.admin.notes,
         error: state => state.error,
         loading: state => state.loading,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
         hasNotes: 'admin/hasNotes',
      })
   },
   created() {
      this.$store.dispatch("admin/getAccessionDetail", this.$route.params.id)
   },
   methods: {
      formattedDate(createdAt) {
         return createdAt.split("T")[0]
      },
      safeCSV(list) {
         if (list) {
            return list.join(", ")
         }
         return "N/A"
      }
   }
}
</script>

<style scoped>
div.admin {
   background: white;
   color: #444;
   position: relative;
   min-width: 1000px;
   padding: 30px 50px 250px 50px;
}
.accession-details {
   padding:0;
   margin:15px 0;
}
h3 {
   margin: 10px 0;
   background: #eee;
   padding: 2px 5px;
   padding: 2px 10px;
   border-radius: 15px;
   color: #666;
}
.info-block {
   margin: 0 0 15px 15px;
}
p {
   margin: 0px 0px 5px 25px;
   font-weight: 100;
   color: #666;
}
a {
   color: cornflowerblue;
   font-weight: 100;
   text-decoration: none;
   font-size:0.9em;
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
div.note {
   padding: 15px;
   margin: 5px;
   border: 1px solid #ccc;
   background: white;
   position: relative;
}
div.note b {
   margin-right: 10px;
}
span.note-time {
   position: absolute;
   right: 15px;
   top: 15px;
   font-size: 0.9em;
}
div.note p {
   margin: 2px;
   border-bottom: 1px solid #ccc;
  margin-bottom: 10px;
  padding-bottom: 5px;
}
</style>