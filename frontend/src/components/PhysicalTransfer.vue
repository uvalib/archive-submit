<template>
   <div class="physical-form">
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
               <label for="box-info">Numer and Size of Boxes</label>
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
               <label>
                  Does this transfer include items on digital media carriers? 
                  <span class="note">(e.g., USB Drives, 5.25-inch floppy disks, CDs)</span>
               </label>
               <label class="pure-radio inline">
                  <input v-model="hasDigital" type="radio" name="has-digital" value="1">
                  Yes
               </label>
               <label class="pure-radio inline">
                  <input v-model="hasDigital" type="radio" name="has-digital" value="0">
                  No
               </label>
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
            <div class="digital-content-questions">
               <div class="pure-u-1-1 gap">
                  <label class="digital-info">
                     Describe Technical Information
                     <span class="note">(e.g., file structure and organization, software that created files, OS, hardware, naming conventions, and original location).</span>
                  </label>
                  <textarea v-model="techInfo" class="digital-info pure-u-1-1" id="tech-description"></textarea>
               </div>
               
               <div class="pure-u-2-5">
                  <label class="digital-info">Do the digital records include any software?</label>
                  <select name="has-software" v-model="hasSoftware">
                     <option value="no">No</option>
                     <option value="yes">Yes</option>
                     <option value="unsure">Unsure</option>
                  </select>
               </div>
               <div class="pure-u-2-5">
                  <label>
                     Please provide an inventory of the items in this transfer
                     <span class="note">(click the inventory button to open the form)</span>
                  </label>
               </div>
               <div class="pure-u-1-5">
                  <span @click="inventoryClicked" class="inventory pure-button pure-button-primary">
                     Inventory Form
                  </span>
                  <span class="note">Inventory Size: {{inventoryCount}}</span>
               </div>
            </div>
         </div>
      </AccordionContent>
      <InventoryForm v-if="showInventory"/>
   </div>
</template>

<script>
import AccordionContent from "@/components/AccordionContent"
import InventoryForm from "@/components/InventoryForm"
import { mapFields } from 'vuex-map-fields'
import { mapState } from "vuex"
import { mapGetters } from 'vuex'

export default {
   components: {
      AccordionContent: AccordionContent,
      InventoryForm: InventoryForm
   },
   watch: {
      hasDigital: function (val) {
         let eles = document.getElementsByClassName("digital-info")
         if (val === "0") {
            this.$store.commit("transfer/clearPhysicalXefrDigitalInfo")
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
      ...mapState({
         mediaCarrierChoices: state => state.transfer.mediaCarrierChoices,
         transferMethods: state => state.transfer.transferMethods,
         physicalRecordTypes: state => state.transfer.physicalRecordTypes,
         showInventory: state => state.transfer.showInventory,
      }),
      ...mapFields([
         'transfer.physical.dateRange',
         'transfer.physical.boxInfo',
         'transfer.physical.selectedTypes',
         'transfer.physical.transferMethod',
         'transfer.physical.hasDigital',
         'transfer.physical.techInfo',
         'transfer.physical.mediaCarriers',
         'transfer.physical.mediaCount',
         'transfer.physical.hasSoftware',
      ]),
      ...mapGetters({
         inventoryCount: 'transfer/inventoryCount',
      })
   },
   created: function() {
      this.$store.dispatch("transfer/getTransferMethods")
      this.$store.dispatch("transfer/getMediaCarriers")
   },
   methods: {
      inventoryClicked() {
         this.$store.commit("transfer/toggleInventory")
      }
   }
};
</script>

<style scoped>
.inventory.pure-button {
   width: 100%;
}
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