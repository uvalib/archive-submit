<template>
   <div id="app">
      <ArchivesHeader/>
      <router-view/>
      <ArchivesFooter/>
   </div>
</template>


<script>
import ArchivesHeader from "@/components/ArchivesHeader";
import ArchivesFooter from "@/components/ArchivesFooter";
import { mapGetters } from "vuex";

export default {
   components: {
      ArchivesHeader,
      ArchivesFooter
   },
   computed: {
      ...mapGetters({
         isAuthenticated: "admin/isAuthenticated"
      })
   },
   created: function() {
      if (this.$route.meta.requiresAuth) {
         if (this.isAuthenticated == false) {
            let authUser = this.$cookies.get("archives_xfer_user");
            if (authUser) {
               authUser.authenticated = true;
               this.$store.commit("setUser", authUser);
               this.$cookies.remove("archives_xfer_user");
               this.$cookies.remove("archives_xfer_settings");
            } else {
               this.$router.push("forbidden");
            }
         }
      }
   }
};
</script>

<style>
html,
body {
   height: 100%;
   margin: 0;
   padding: 0;
   background: #f8f2e3;
}
div.main-content {
   padding-bottom: 210px;
}
#app {
   font-family: "Avenir", Helvetica, Arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   min-height: 100%;
   position: relative;
}
#app h2 {
   margin: 0;
   color: #51822f;
   font-weight: bold;
   position: relative;
   font-size: 24px;
}
#customdropzone {
   padding: 5px;
   height: 180px;
}
#customdropzone .dz-preview {
   display: inline-block;
   background: transparent;
   margin: 5px;
}
#customdropzone .dz-preview .dz-details {
   background-color: rgba(64, 64, 64, 0.5);
   transition: none;
}
#customdropzone .dz-preview .dz-image {
   width: 150px;
   height: 150px;
}
#app a.dz-remove {
   color: white;
}
#app div.content {
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
   color: #444;
}
</style>

