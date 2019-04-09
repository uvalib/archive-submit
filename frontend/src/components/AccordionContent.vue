<template>
   <div class="accordion">
      <h4 @click="accordionClicked">
         {{title}}
         <span class="note">{{subtitle}}</span>
         <i class="accordion-icon fas fa-angle-down" v-bind:style="{ transform: rotation }"></i>
      </h4>
       <transition name="accordion"
         v-on:before-enter="beforeEnter" v-on:enter="enter"
         v-on:before-leave="beforeLeave" v-on:leave="leave">
         <div class="accordion-content"  v-show="expanded">
            <slot></slot>
         </div>
      </transition>
   </div>
</template>

<script>
export default {
   props: {
      title: String,
      subtitle: String
   },
   data: function() {
      return {
         expanded: false
      };
   },
   computed: {
      rotation() {
         if (this.expanded) {
            return "rotate(180deg)"
         }
         return "rotate(0deg)"
      }
   },
   methods: {
      accordionClicked() {
         this.expanded = !this.expanded
      },
      beforeEnter: function(el) {
         el.style.height = '0'
      },
      enter: function(el) {
         el.style.height = el.scrollHeight + 'px'
      },
      beforeLeave: function(el) {
         el.style.height = el.scrollHeight + 'px'
      },
      leave: function(el) {
         el.style.height = '0'
      }
   }
};
</script>

<style scoped>
.accordion-content {
   overflow: hidden;
   transition: 500ms ease-out;
}
h4 {
   padding: 5px;
   position: relative;
   cursor: pointer;
   border-bottom: 1px solid #666;
   margin: 10px 0;
}
h4:hover {
   color: #444;
   border-bottom: 1px solid #444;
}
.note {
   font-weight: 100;
}
h4 .accordion-icon {
   font-size: 1.25em;
   position: absolute;
   right: 5px;
   cursor: pointer;
   transform: rotate(0deg);
   transition-duration: 0.7s;
}
</style>