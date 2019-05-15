<template>
   <AccordionContent title="Notes">
      <button>Add Note</button>
      <div class="note" v-for="(note, idx) in notes" :key="idx">
         <p><b>{{note.title}}</b></p>
         <span class="note-time"><b>{{note.userName}}</b>{{formattedDate(note.createdAt)}}</span>
         <div>{{note.note}}</div>
      </div>
   </AccordionContent>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
import AccordionContent from '@/components/AccordionContent'
export default {
   components: {
      AccordionContent: AccordionContent,
   },
   computed: {
      ...mapState({
         notes: state=>state.admin.notes,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
         hasNotes: 'admin/hasNotes',
      })
   },
   methods: {
      formattedDate(createdAt) {
         return createdAt.split("T")[0]
      },
   }
};
</script>

<style scoped>
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