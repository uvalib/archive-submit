<template>
   <transition name="modal">
      <div class="modal-mask">
         <div class="modal-container">
            <h3>Physical Records Inventory</h3>
            <div class="content">
               <table class="pure-table  pure-table-bordered"> 
                  <thead>
                     <tr>
                        <th>Box Number</th>
                        <th>Records Group #</th>
                        <th style="width:175px">Box Title</th>
                        <th style="width:400px">Contents Description</th>
                        <th>Dates</th>
                        <th/>
                     </tr>
                  </thead>
                  <tr v-for="(item, idx) in inventory" :key="idx">
                     <template v-if="editIdx != idx">
                        <td>{{item.boxNum}}</td>
                        <td>{{item.recordGroup}}</td>
                        <td>{{item.title}}</td>
                        <td>{{item.description}}</td>
                        <td>{{item.dates}}</td>
                        <td style="text-align:center; width: 40px; padding:8px 4px;">
                           <i @click="editClicked" class="fas fa-edit" :data-idx="idx"></i>
                           <i @click="removeClicked" class="right fas fa-trash-alt" :data-idx="idx"></i>
                        </td> 
                     </template>
                     <template v-else>
                        <td class="edit" ><input id="first" type="text" v-model="editItem.boxNum"/></td>
                        <td class="edit" ><input type="text" v-model="editItem.recordGroup"/></td>
                        <td class="edit" ><input type="text" v-model="editItem.title"/></td>
                        <td class="edit" ><textarea v-model="editItem.description"></textarea></td>
                        <td class="edit" ><input type="text" v-model="editItem.dates"/></td>
                        <td style="text-align:center; width: 40px; padding:8px 4px;">
                           <i @click="okClicked" class="fas fa-check-circle"></i>
                           <i @click="cancelClicked" class="right fas fa-window-close"></i>
                        </td> 
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
         editItem: null
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
         this.editItem = this.$store.getters.inventoryItem(idx)
         setTimeout( function() {
            document.getElementById("first").focus()
         }, 100)
      },
      cancelClicked() {
         this.editIdx = -1
         this.editItem = null
      },
      okClicked() {
          this.$store.commit("updateInventory", {idx:this.editIdx, item:this.editItem} )
          this.editIdx = -1
          this.editItem = null
      },
   }
};
</script>

<style scoped>
i.fas {
   cursor: pointer;
   opacity: 0.5;
}
i.fas.right {
   margin-left: 5px;
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
.content table td input,  .content table td textarea {
   margin: 0;
   width: 100%;
   border-radius: 0;
   box-shadow: none;
   border: none;
   outline:none;
   background: aliceblue;
}
.content table td.edit {
   padding:0px;
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
