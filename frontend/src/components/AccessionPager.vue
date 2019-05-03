<template>
   <div class="pager">
      <span class="total">{{ total }}&nbsp;accessions</span>
      <i class="button fas fa-angle-double-left"></i><i class="button fas fa-angle-left"></i>
      <span class="curr">{{ page }} of {{ lastPage }}</span>
      <i class="button fas fa-angle-right"></i><i class="button fas fa-angle-double-right"></i>
   </div>
</template>

<script>
import { mapState } from "vuex";
export default {
   computed: {
      ...mapState({
         total: state => state.admin.totalAccessions,
         page: state => state.admin.page,
         pageSize: state => state.admin.pageSize,
      }),
      lastPage() {
         return  Math.floor( this.total / this.pageSize)+1
      }
   },
   methods: {
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
span.curr {
   display: inline-block;
   border: 1px solid #ccc;
   padding: 1px 10px;
   height: 16px;
   margin: 0 2px;
}
</style>