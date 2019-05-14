<template>
   <div class="header">
      <router-link :to="{ name: 'home'}">
         <span class="site-name small">Albert and Shelly Small Special Collections Library</span>
         <span class="site-name">Records Transfer Form</span>
      </router-link>
      <a target="_blank" href="https://library.virginia.edu"><img class="uva-library" src="../assets/uvalogo.png"/></a>
   </div>
</template>

<script>
export default {
}
</script>

<style scoped>
div.header {
   background-color: #20406e;
   color: white;
   padding:15px;
   border-bottom: 1px solid #5e7799;
   border-top: 1px solid #5e7799;
   text-align: left;
}
a {
  color: white;
  text-decoration: none;
}
div.header span.site-name {
   margin: 0;
   font-family: "chaparral-pro", Georgia, "Times New Roman", Times, serif;
   font-size: 34px;
   position: relative;
}
div.header span.site-name.small {
  color: #bfc9d5;
  font-size: 19px;
  line-height: 19px;
  margin-bottom: 0;
  display: block;
}
img.uva-library {
  position: absolute;
  top: 10px;
  right: 10px;
  max-height: 70px;
}
.subtitle {
   display: inline-block;
   margin: 5px 0 0 15px;
}
</style>
