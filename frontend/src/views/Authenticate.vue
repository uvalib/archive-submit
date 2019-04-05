<template>
   <div class="authenticate">
      <h2>University Archives Records Transfer Form</h2>
      <div>
         <h3>Authenticating...</h3>
      </div>
   </div>
</template>

<script>
import axios from "axios";
export default {
  name: 'authenticate',
  created: function () {
     axios
         .post("/auth/verify")
         .then(response => {
            this.$store.commit("setUser", response.data);
            this.$router.push("submit");
         })
         .catch((/*error*/) => {
            this.$router.push("forbidden");
         });
  }
}
</script>

<style scoped>
div.authenticate {
   position: relative;
   min-width: 1000px;
   padding: 30px 50px 250px 50px;
   min-height: 600px;
   background: white;
   width: 75%;
   margin: 0 auto;
   border-right: 1px solid #dfdacc;
   border-left: 1px solid #dfdacc;
   font-weight: 400;
   color: #333;
}
</style>