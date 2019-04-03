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
         <span class="note">(creating office, department, center, institute, or administrator(s)):</span>
         <input class="pure-u-1-1" id="activities" type="text"/>
      </div>
   </div>
</template>

<script>
  export default {
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
</style>