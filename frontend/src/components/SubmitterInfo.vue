<template>
   <div class="submitter">
      <div class="spacing">
         <div class="pure-u-1-3">
            <label for="lname">Last Name</label>
            <input id="lname" class="pure-u-23-24" type="text" required :value="user.lastName">
         </div>
         <div class="pure-u-1-3">
            <label for="fname">First Name</label>
            <input id="fname" class="pure-u-23-24" type="text" required :value="user.firstName">
         </div>
         <div class="pure-u-1-3">
            <label for="email">Email Address</label>
            <input id="email" class="pure-u-23-24" type="email" required :value="user.email">
         </div>
      </div>
      <div class="spacing">
         <div class="pure-u-1-3">
            <label for="title">Title</label>
            <input class="pure-u-23-24" id="title" type="text" required :value="user.title">
         </div>
         <div class="pure-u-1-3">
            <label for="affiliation">University Affiliation</label>
            <input class="pure-u-23-24" id="affiliation" type="text" required :value="user.affiliation">
         </div>
         <div class="pure-u-1-3">
            <label for="phone">Phone Number</label>
            <input class="pure-u-23-24" id="phone" type="tel" required :value="user.phone">
         </div>
      </div>
   </div>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
   computed: {
      ...mapGetters([
         'user'
      ])
   }
}
</script>

<style scoped>
div.submitter {
   width:100%;
}
div.submitter label {
   color:#666;
   display: block;
}
div.spacing {
   margin-bottom: 10px;
}
</style>
