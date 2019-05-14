<template>
   <div class="pager">
      <span v-if="filteredTotal==0" class="total">{{ total }}&nbsp;accessions</span>
      <span v-else class="total">{{filteredTotal}} filtered from {{ total }}&nbsp;accessions</span>
      <i @click="firstPageClicked" v-bind:class="{disabled: isFirstPage}" class="button fas fa-angle-double-left"></i>
      <i @click="prevPageClicked" v-bind:class="{disabled: isFirstPage}" class="button fas fa-angle-left"></i>
      <span class="curr">{{ page }} of {{ lastPage }}</span>
      <i @click="nextPageClicked" v-bind:class="{disabled: isLastPage}" class="button fas fa-angle-right"></i>
      <i @click="lastPageClicked" v-bind:class="{disabled: isLastPage}" class="button fas fa-angle-double-right"></i>
   </div>
</template>

<script>
import { mapState } from "vuex";
export default {
   computed: {
      ...mapState({
         total: state => state.admin.totalAccessions,
         filteredTotal: state => state.admin.filteredTotal,
         page: state => state.admin.page,
         pageSize: state => state.admin.pageSize,
         queryStr: state => state.admin.queryStr,
      }),
      lastPage() {
         if ( this.filteredTotal > 0) {
            return  Math.floor( this.filteredTotal / this.pageSize)+1
         }
         return  Math.floor( this.total / this.pageSize)+1
      },
      isFirstPage() {
         return this.page === 1
      },
      isLastPage() {
         return this.page === this.lastPage
      }
   },
   methods: {
      firstPageClicked() {
         this.$store.dispatch("admin/firstPage")
      },
      lastPageClicked() {
         this.$store.dispatch("admin/lastPage")
      },
      nextPageClicked() {
         this.$store.dispatch("admin/nextPage")
      },
      prevPageClicked() {
         this.$store.dispatch("admin/prevPage")
      }
   }
};
</script>

<style scoped>
.pager {
   position: absolute;
   right:0;
   top:0;
   font-size: 12px;
}
.total {
   margin-right: 10px;
   display: inline-block;
}
.button.fas {
   width: 20px;
   height: 12px;
   padding: 4px 0;
   text-align: center;
   margin: 1px;
   background: #0078e7;
   color: white;
}
.button.fas:hover {
   cursor:  pointer !important;
   background: #1088f7;
}
.button.fas.disabled {
   cursor:default;
   opacity: 0.5;
}
.button.fas.disabled:hover {
   cursor:default;
   opacity: 0.5;
}
span.curr {
   display: inline-block;
   border: 1px solid #ccc;
   padding: 1px 10px;
   height: 16px;
   margin: 0 2px;
}
</style>