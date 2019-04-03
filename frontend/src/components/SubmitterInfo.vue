<template>
   <div class="submitter">
      <div class="spacing">
         <div class="pure-u-1-3">
            <label for="lname">Last Name</label>
            <input id="lname" class="pure-u-23-24" type="text" required>
         </div>
         <div class="pure-u-1-3">
            <label for="fname">First Name</label>
            <input id="fname" class="pure-u-23-24" type="text" required>
         </div>
         <div class="pure-u-1-3">
            <label for="email">Email Address</label>
            <input id="email" class="pure-u-23-24" type="email" required>
         </div>
      </div>
      <div class="spacing">
         <div class="pure-u-1-3">
            <label for="title">Title</label>
            <input class="pure-u-23-24" id="title" type="text">
         </div>
         <div class="pure-u-1-3">
            <label for="affiliation">University Affiliation</label>
            <input class="pure-u-23-24" id="affiliation" type="text">
         </div>
         <div class="pure-u-1-3">
            <label for="phone">Phone Number</label>
            <input class="pure-u-23-24" id="phone" type="tel">
         </div>
      </div>
   </div>
</template>

<script>
  export default {
  }
</script>

<style scoped>
div.submitter {
   width:100%;
}
div.spacing {
   margin-bottom: 10px;
}
</style>
