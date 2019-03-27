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
    }
  }
</script>

<style>
  html,body{
      height: 100%;
      margin: 0;
      padding: 0;
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
</style>

