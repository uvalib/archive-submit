<template>
   <transition name="modal">
      <div class="modal-mask">
         <div class="modal-container">
            <h3>Physical Records Inventory</h3>
            <div class="content">
               <table class="pure-table  pure-table-bordered"> 
                  <thead>
                     <tr>
                        <th/>
                        <th>Box Number</th>
                        <th>Records Group #</th>
                        <th>Box Title</th>
                        <th>Contents Description</th>
                        <th>Dates</th>
                     </tr>
                  </thead>
                  <tr v-for="(item, idx) in inventory" :key="idx">
                     <template v-if="editIdx != idx">
                        <td style="text-align:center">
                           <i @click="editClicked" class="fas fa-edit" :data-idx="idx"></i>
                           <i @click="removeClicked" class="fas fa-trash-alt" :data-idx="idx"></i>
                        </td> 
                        <td>readonly</td>
                        <td></td>
                        <td></td>
                        <td></td>
                        <td></td>
                     </template>
                     <template v-else>
                        <td style="text-align:center">
                           <i @click="okClicked" class="fas fa-check-circle"></i>
                           <i @click="cancelClicked" class="fas fa-window-close"></i>
                        </td> 
                        <td>EDIT</td>
                        <td></td>
                        <td></td>
                        <td></td>
                        <td></td>
                     </template>
                  </tr>
               </table>
            </div>
            <div class="controls">
               <span @click="addClicked" class="pure-button pure-button-primary">Add Item</span>
               <span @click="clearClicked" class="pure-button pure-button-primary">Clear Form</span>
               <span @click="closeClicked" class="pure-button pure-button-primary">OK</span>
            </div>
         </div>
      </div>
   </transition>
</template>

<script>
export default {
   components: {
   },
   data: function () {
      return {
         editIdx: -1,
         origItem: null
      }
   },
   computed: {
      inventory() {
         return this.$store.getters.inventory
      },
   },
   methods: {
      closeClicked() {
         this.$store.commit("toggleInventory")
      },
      clearClicked() {
         this.$store.commit("clearInventory")
      },
      addClicked() {
         this.$store.commit("addInventory")
      },
      removeClicked(event) {
         let idx = event.target.dataset.idx
         this.$store.commit("deleteInventory", idx)
      },
      editClicked(event) {
         let idx = event.target.dataset.idx
         this.editIdx = idx
         this.origItem = this.$store.getters.inventoryItem(idx)
         alert(this.origItem)
      },
      cancelClicked() {
         this.editIdx = -1
      },
      okClicked() {
      },
   }
};
</script>

<style scoped>
i.fas {
   cursor: pointer;
   opacity: 0.5;
   margin-left:10px;
}
i.fas:hover {
   opacity: 1;
}
.modal-mask {
   position: fixed;
   z-index: 9998;
   top: 0;
   left: 0;
   width: 100%;
   height: 100%;
   background-color: rgba(0, 0, 0, 0.5);
   display: table;
   transition: opacity 0.2s ease;
}
span.pure-button  {
   margin-left:10px;
}
table {
   width:100%;
}
.content {
   padding: 10px;
}
h3 {
   padding: 2px 8px;
   background-color: #20406e;
   border-bottom: 1px solid #5e7799;
   color:white;
   margin:0;
}
.controls {
   text-align: right;
   padding: 0px 10px 10px 10px;
}
.modal-container {
   font-size: 0.9em;
   width: 90%;
   margin: 2% auto;
   padding: 0px;
   background-color: #fff;
   box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
}
</style>
