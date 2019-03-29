<template>
   <div class="home">
      <h2>
         <span>Archives Records Transfer System</span>
         <span class="admin">
            <router-link to="/admin">
               <button class="pure-button pure-button-primary">Admin Access</button>
            </router-link>
         </span>
      </h2>
      <div class="access-type">
         <h3>Access the submission form</h3>
         <label>
            <input type="radio" name="is_uva" id="uva" value="no" checked="checked">
            I am UVA faculty, staff, or student. (You will be asked to verify your identity using 
            <a href="http://itc.virginia.edu/netbadge/">NetBadge</a>.)
         </label>
         <label>
            <input type="radio" name="is_uva" id="guest" value="yes" checked="checked">
            I am not affiliated with UVA.
         </label>
         <router-link to="/submit">
            <button class="pure-button pure-button-primary">Continue</button>
         </router-link>
      </div>
   </div>
</template>

<script>
export default {
  name: 'home',
}
</script>

<style scoped>
a {
   color: cornflowerblue;
   font-weight: 500;
   text-decoration: none;
}
span.admin {
   font-size: 12px;
   position: absolute;
}
a:hover {
 text-decoration: underline;  
}
div.home {
  padding: 30px 50px;;
  min-height: 600px;
  background: white;
  width: 75%;
  margin: 0 auto;
  border-right: 1px solid #dfdacc;
  border-left: 1px solid #dfdacc;
  font-weight: 400;
  color: #333;
}
h2 {
  margin: 0;
  color:#51822F;
  font-weight: bold;
  position: relative;
}
h3 {
   color:#EB5F0C;
   border-bottom: 1px dashed #EB5F0C;
   padding-bottom: 5px;
   font-size: 1.3em;
}
.access-type input {
  margin-right: 10px;
}
.access-type label {
   display: block;
   margin: 0 0 10px 25px;
}
</style>