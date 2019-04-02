<template>
  <div id="app">
    <ArchivesHeader/>
    <router-view/>
    <ArchivesFooter/>
  </div>
</template>


<script>
  import ArchivesHeader from '@/components/ArchivesHeader'
  import ArchivesFooter from '@/components/ArchivesFooter'

  export default {
    components: {
      ArchivesHeader,
      ArchivesFooter
    },
    data: function () {
      return {
        // authorizing: true,
        authorizing: false,
      }
    },
    created: function () {
      // console.log("CREATED")
      // auth checks here see also apollo backend handler for authenticate and the auth_middleware
       // if (this.$route.meta.requiresAuth ) {
    //   axios.get("/authenticate").then((response)  =>  {
    //     localStorage.setItem("user", response.data.firstName+" "+response.data.lastName)
    //     this.authorizing = false
    //   }).catch(() => {
    //     this.authorizing = false
    //     localStorage.removeItem("user")
    //     this.$router.push({ path: 'unauthorized' })
    //   })
    // } else {
    //   this.authorizing = false
    // }
    }
  }
</script>

<style>
  html,body{
      height: 100%;
      margin: 0;
      padding: 0;
      background: #f8f2e3;
  }
  div.main-content {
    padding-bottom: 210px;
  }
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    min-height: 100%;
    position:relative;
  }
  #app h2 {
    margin: 0;
    color:#51822F;
    font-weight: bold;
    position: relative;
    font-size: 24px;
  }
</style>

