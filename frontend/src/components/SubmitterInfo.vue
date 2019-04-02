<template>
   <div class="submitter">
      <div class="spacing">
         <div class="pure-u-1-3">
            <label for="lname">Last Name</label>
            <input id="lname" class="pure-u-23-24" type="text" placeholder="Last Name" required>
         </div>
         <div class="pure-u-1-3">
            <label for="fname">First Name</label>
            <input id="fname" class="pure-u-23-24" type="text" placeholder="First Name" required>
         </div>
         <div class="pure-u-1-3">
            <label for="email">Email Address</label>
            <input id="email" class="pure-u-23-24" type="email" required placeholder="Email Address" >
         </div>
      </div>
      <div class="spacing">
         <div class="pure-u-1-3">
            <label for="title">Title</label>
            <input class="pure-u-23-24" id="title" type="text" placeholder="Title">
         </div>
         <div class="pure-u-1-3">
            <label for="affiliation">University Affiliation</label>
            <input class="pure-u-23-24" id="affiliation" type="text" placeholder="Affiliation">
         </div>
         <div class="pure-u-1-3">
            <label for="phone">Phone Number</label>
            <input class="pure-u-23-24" id="phone" type="tel" placeholder="Phone Number">
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
