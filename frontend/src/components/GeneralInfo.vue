<template>
   <div class="general-info">
      <div class="pure-u-1-1 bottom-pad ">
         <label for="description">Summary Description</label>
         <span class="note">(e.g., Title, Types of Materials, Nature, Item Relationships, Duplicated/Missing Materials, Personally Identifiable Information)</span>
         <textarea class="pure-u-1-1" id="description"></textarea>
      </div>
      <div class="pure-u-1-1 bottom-pad ">
         <label for="activities">What were the activities that led to the creation of the records?</label>
         <textarea class="pure-u-1-1" id="activities"></textarea>
      </div>
      <div class="pure-u-1-1 bottom-pad ">
         <label for="creator">Creator of the records</label>
         <span class="note">(creating office, department, center, institute, or administrator(s))</span>
         <input class="pure-u-1-1" id="activities" type="text"/>
      </div>
      <div class="pure-u-1-2 bottom-pad ">
         <label for="creator">Genres of Materials</label>
         <span class="note">(check all that apply)</span>
         <div class="choices">
            <span v-for="genre in genres" :key="genre.id">
               <label class="pure-checkbox inline">
                  <input type="checkbox" name="genre" :value="genre.id">
                  {{ genre.name }}
               </label>
            </span>
         </div>
      </div>
      <div class="pure-u-1-2 bottom-pad ">
         <label for="creator">Is this a new accession or an accrual to an existing record group? </label>
         <div class="choices">
            <label class="pure-radio">
               <input type="radio" name="accession-type" value="new" checked>
               New accession
            </label>
            <label class="pure-radio">
               <input type="radio" name="accession-type" value="add">
               Addition to existing record group    
            </label>
            <label class="pure-radio">
               <input type="radio" name="accession-type" value="unsure">
               Unsure
            </label>
         </div>
      </div>
   </div>
</template>

<script>
  export default {
   computed: {
      genres: function() {
         return this.$store.getters.genres
      }
   }
  }
</script>

<style scoped>
div.general-info {
   width:100%;
}
div.bottom-pad {
  margin-bottom: 10px;
}
span.note {
  color: #999;
  font-size:0.85em;
}
div.choices {
   padding: 5px 0;
}
label.pure-checkbox.inline {
   display: inline-block;
   margin: 5px 15px 5px 0;
}
</style>