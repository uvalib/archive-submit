<template>
   <div class="header">
      <div class="site-link">
         <router-link to="/">
            <span class="site-name small">Albert and Shelly Small Special Collections Library</span>
            <span class="site-name">Records Transfer Form</span>
         </router-link>
      </div>
      <div class="library-link">
         <a target="_blank" href="https://library.virginia.edu">
            <img class="uva-library" src="../assets/uvalogo.png"/>
         </a>
      </div>
      <div style="clear:both"></div>
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
   padding:10px;
   border-bottom: 1px solid #5e7799;
   border-top: 1px solid #5e7799;
   text-align: left;
   position: relative;
   box-sizing: border-box;
}
a {
  color: white;
  text-decoration: none;
}
div.header span.site-name {
   margin: 0;
   font-family: "chaparral-pro", Georgia, "Times New Roman", Times, serif;
   font-size: 28px;
   position: relative;
}
div.site-link {
   float:left;
}
div.header span.site-name.small {
  color: #bfc9d5;
  font-size: 19px;
  line-height: 19px;
  margin-bottom: 0;
  display: block;
}
img.uva-library {
  max-height: 50px;
}
div.library-link {
   float:right;
}
.subtitle {
   display: inline-block;
   margin: 5px 0 0 15px;
}
</style>
