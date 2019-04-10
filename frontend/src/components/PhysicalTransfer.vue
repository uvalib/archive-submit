<template>
   <AccordionContent
      title="Physical Records Transfer"
      subtitle="(Paper-based/Analog Records, Media Carriers)"
   >
      <div class="data-form">
         <div class="pure-u-1-2">
            <label for="phys-date-range">Date Range of Records</label>
            <input id="phys-date-range" class="pure-u-23-24" type="text">
         </div>
         <div class="pure-u-1-2">
            <label for="box-info">Numer and size of Boxes</label>
            <input id="box-info" class="pure-u-23-24" type="text">
         </div>
         <div class="pure-u-1-1 gap">
            <label>
               Record Types
               <span class="note">(check all that apply)</span>
            </label>
            <div class="choices">
               <span v-for="rt in physicalRecordTypes" :key="rt.id">
                  <label class="pure-checkbox inline">
                     <input type="checkbox" name="phys-record-type" :value="rt.id">
                     {{ rt.name }}
                     <span class="note">{{rt.description}}</span>
                  </label>
               </span>
            </div>
         </div>
         <div class="pure-u-1-1 gap">
            <label>Transfer Method</label>
            <div class="choices">
               <span v-for="m in transferMethods" :key="m.id">
                  <label class="pure-checkbox inline">
                     <input type="checkbox" name="transfer-method" :value="m.id">
                     {{ m.name }}
                  </label>
               </span>
            </div>
         </div>
         <div class="pure-u-1-1">
            <label>Does this transfer include items on digital media carriers?</label>
            <label class="pure-radio inline">
               <input v-model="digitalContent" type="radio" name="has-digital" value="yes" checked>
               Yes
            </label>
            <label class="pure-radio inline">
               <input v-model="digitalContent" type="radio" name="has-digital" value="no">
               No
            </label>
         </div>
         <div v-if="digitalContent === 'yes'" class="digital-content-questions">
            <div class="pure-u-1-1 gap">
               <label>
                  Describe Technical Information
                  <span class="note">(e.g., file structure and organization, software that created files, OS, hardware, naming conventions, and original location).</span>
               </label>
               <textarea class="pure-u-1-1" id="tech-description"></textarea>
            </div>
            
            <div class="pure-u-1-1">
               <label>
                  Media carriers
                  <span class="note">(check all that apply)</span>
               </label>
               <div class="choices">
                  <span v-for="mc in mediaCarriers" :key="mc.id">
                     <label class="pure-checkbox inline">
                        <input type="checkbox" name="media-carrier" :value="mc.id">
                        {{ mc.name }}
                     </label>
                  </span>
               </div>
            </div>
            <div class="pure-u-1-1 gap">
               <label>
                  Please estimate the number of each media format and their types 
                  <span class="note">(e.g., 5 1/4" disk, 3 1/2" disk, CD-ROM, thumb drive, etc.)</span>
               </label>
               <textarea class="pure-u-1-1" id="media-count"></textarea>
            </div>
            <div class="pure-u-1-1">
               <label>Do the digital records include any software?</label>
               <label class="pure-radio inline">
                  <input type="radio" name="has-software" value="yes" checked>
                  Yes
               </label>
               <label class="pure-radio inline">
                  <input type="radio" name="has-software" value="no">
                  No
               </label>
            </div>
         </div>
      </div>
   </AccordionContent>
</template>

<script>
import AccordionContent from "@/components/AccordionContent";
export default {
   components: {
      AccordionContent: AccordionContent
   },
   data: function() {
      return {
         digitalContent: "yes"
      };
   },
   computed: {
      physicalRecordTypes() {
         return this.$store.getters.physicalRecordTypes;
      },
      transferMethods() {
         return this.$store.getters.transferMethods;
      },
      mediaCarriers() {
         return this.$store.getters.mediaCarriers;
      }
   },
   created: function() {
      this.$store.dispatch("getTransferMethods");
      this.$store.dispatch("getMediaCarriers");
   }
};
</script>

<style scoped>
.note {
   color: #999;
  font-size: 0.85em;
}
div.data-form {
   padding: 10px;
}
div.gap {
   margin: 10px 0;
}
.inline {
   display: inline-block;
   margin-left: 15px;
}
</style>