<template>
   <div>
      <AccordionContent  v-if="hasNotes" title="Notes" :watched="notes">
         <div class="note" v-for="(note, idx) in notes" :key="idx">
            <p><b>{{note.title}}</b></p>
            <span class="note-time"><b>{{note.userName}}</b>{{formattedDate(note.createdAt)}}</span>
            <div>{{note.note}}</div>
         </div>
      </AccordionContent>
      <div class="note-actions">
         <template v-if="addingNote">
            <template v-if="working">
               <span class="working">Adding note...</span>
            </template>
            <template v-else>
               <span @click="cancelNote" class="cancel-note pure-button pure-button-primary">Cancel</span>
               <span @click="addNote" class="add-note pure-button pure-button-primary">Submit</span>
            </template>
         </template>
         <template v-else>
            <span @click="showAddNote" class="show-add pure-button pure-button-primary">Add Note</span>
         </template>
      </div>
      <div v-if="addingNote" class="add-form pure-form-aligned">
         <div class="pure-control-group">
            <label for="title">Title</label>
            <input id="title" type="text" v-model="newTitle">
        </div>
        <div class="pure-control-group">
            <label for="title">Note</label>
            <textarea rows="5" id="note" v-model="newNote"></textarea>
        </div>
        <p class="error">{{error}}</p>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
import AccordionContent from '@/components/AccordionContent'
export default {
   data: function() {
      return {
         newTitle: "",
         newNote: "",
      };
   },
   components: {
      AccordionContent: AccordionContent,
   },
   computed: {
      ...mapState({
         addingNote: state=>state.admin.addingNote,
         notes: state=>state.admin.notes,
         error: state=>state.error,
         working: state=>state.admin.working,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
         hasNotes: 'admin/hasNotes',
      })
   },
   methods: {
      showAddNote() {
         this.newTitle = ""
         this.newNote = ""
         this.$store.commit("admin/setAddingNote", true)
      },
      addNote(event) {
         event.stopPropagation()
         if (this.newNote.length == 0 || this.newTitle.length == 0) {
            this.$store.commit("setError", "Title and note text are required")
            return;
         }
         this.$store.dispatch("admin/addNote", {title: this.newTitle, note: this.newNote})
      },
      cancelNote() {
        this.$store.commit("admin/setAddingNote", false)
      },
      formattedDate(createdAt) {
         return createdAt.split("T")[0]
      },
   }
};
</script>

<style scoped>
.working {
   color: #999;
   font-style: italic;
   font-weight: bold;
}
p.error {
   color: firebrick;
   text-align: center;  
}
span.pure-button.pure-button-primary {
   font-size: 0.8em;
}
span.show-add {
     margin-right: 5px;
}
.cancel-note {
   margin-right: 10px;
}
.add-form {
   margin-top: 15px;
}
#title {
   width: 65%;
}
#note {
   border-color: #ccc;
   width: 65%;
}
.note-actions {
   text-align: right;
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