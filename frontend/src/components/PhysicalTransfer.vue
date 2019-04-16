<template>
   <AccordionContent
      title="Physical Records Transfer"
      subtitle="(Paper-based/Analog Records, Media Carriers)"
   >
      <div class="data-form">
         <div class="pure-u-1-2">
            <label for="phys-date-range">Date Range of Records</label>
            <input id="phys-date-range" class="pure-u-23-24" type="text" v-model="dateRange">
         </div>
         <div class="pure-u-1-2">
            <label for="box-info">Numer and size of Boxes</label>
            <input id="box-info" class="pure-u-23-24" type="text" v-model="boxInfo">
         </div>
         <div class="pure-u-1-1 gap">
            <label>
               Record Types
               <span class="note">(check all that apply)</span>
            </label>
            <div class="choices">
               <span v-for="rt in physicalRecordTypes" :key="rt.id">
                  <label class="pure-checkbox inline">
                     <input type="checkbox" name="phys-record-type" :value="rt.id" v-model="selectedTypes">
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
                  <label class="pure-radio inline">
                     <input type="radio" name="transfer-method" :value="m.id" v-model="transferMethod">
                     {{ m.name }}
                  </label>
               </span>
            </div>
         </div>
         <div class="pure-u-1-1">
            <label>Does this transfer include items on digital media carriers?</label>
            <label class="pure-radio inline">
               <input v-model="hasDigital" type="radio" name="has-digital" value="1">
               Yes
            </label>
            <label class="pure-radio inline">
               <input v-model="hasDigital" type="radio" name="has-digital" value="0">
               No
            </label>
         </div>
         <div class="digital-content-questions">
            <div class="pure-u-1-1 gap">
               <label class="digital-info">
                  Describe Technical Information
                  <span class="note">(e.g., file structure and organization, software that created files, OS, hardware, naming conventions, and original location).</span>
               </label>
               <textarea v-model="techInfo" class="digital-info pure-u-1-1" id="tech-description"></textarea>
            </div>
            
            <div class="pure-u-1-1">
               <label class="digital-info">
                  Media carriers
                  <span class="note">(check all that apply)</span>
               </label>
               <div class="choices">
                  <span v-for="mc in mediaCarrierChoices" :key="mc.id">
                     <label class="digital-info pure-checkbox inline">
                        <input type="checkbox" class="digital-info" name="media-carrier" :value="mc.id" v-model="mediaCarriers">
                        {{ mc.name }}
                     </label>
                  </span>
               </div>
            </div>
            <div class="pure-u-1-1 gap">
               <label class="digital-info">
                  Please estimate the number of each media format and their types 
                  <span class="note">(e.g., 5 1/4" disk, 3 1/2" disk, CD-ROM, thumb drive, etc.)</span>
               </label>
               <textarea v-model="mediaCount" class="digital-info pure-u-1-1" id="media-count"></textarea>
            </div>
            <div class="pure-u-1-1">
               <label class="digital-info">Do the digital records include any software?</label>
               <label class="digital-info pure-radio inline">
                  <input type="radio" class="digital-info" name="has-software" value="1" v-model="hasSoftware">
                  Yes
               </label>
               <label class="digital-info pure-radio inline">
                  <input type="radio" class="digital-info" name="has-software" value="0" checked v-model="hasSoftware">
                  No
               </label>
            </div>
         </div>
      </div>
   </AccordionContent>
</template>

<script>
import AccordionContent from "@/components/AccordionContent"
import { mapFields } from 'vuex-map-fields'

export default {
   components: {
      AccordionContent: AccordionContent
   },
   watch: {
      hasDigital: function (val) {
         let eles = document.getElementsByClassName("digital-info")
         if (val === "0") {
            this.$store.commit("clearPhysicalXefrDigitalInfo")
         } 
         Array.from(eles).forEach( function(ele) {
            ele.readOnly = (val === "0")
            ele.disabled = (val === "0")
            if (val == "0") {
               ele.classList.add("ghosted")
            } else {
                ele.classList.remove("ghosted")
            }
         })
      }
   },
   computed: {
      physicalRecordTypes() {
         return this.$store.getters.physicalRecordTypes
      },
      transferMethods() {
         return this.$store.getters.transferMethods
      },
      mediaCarrierChoices() {
         return this.$store.getters.mediaCarrierChoices
      },
      ...mapFields([
         'physical.dateRange',
         'physical.boxInfo',
         'physical.selectedTypes',
         'physical.transferMethod',
         'physical.hasDigital',
         'physical.techInfo',
         'physical.mediaCarriers',
         'physical.mediaCount',
         'physical.hasSoftware',
      ])
   },
   created: function() {
      this.$store.dispatch("getTransferMethods")
      this.$store.dispatch("getMediaCarriers")
   }
};
</script>

<style scoped>
.digital-info.ghosted {
   opacity: 0.5;
}
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