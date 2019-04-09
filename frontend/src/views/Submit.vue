<template>
  <div class="submit">
    <h2>University Archives Records Transfer Form</h2>
    <div class="contact-small">
      <p>University Archives</p>
      <p>Albert & Shirley Small Special Collections Library, P.O. Box 400110, Charlottesville, VA, 22904-4110</p>
      <p>Contact: Bethany Anderson, University Archivist, Phone: (434) 982-2980, Email: <a href="mailto:bga3d@virginia.edu">bga3d@virginia.edu</a></p>
    </div>
    <form class="pure-form pure-form-stacked">
      <fieldset>
        <div class="pure-g">
          <SubmitterInfo/>
          <GeneralInfo/>
        </div>
        <template v-if="showDigital">
          <DigitalTransfer/>
        </template>
        <template v-if="showPhysical">
          <PhysicalTransfer/>
        </template>
      </fieldset>
    </form>
  </div>
</template>

<script>
import SubmitterInfo from '@/components/SubmitterInfo'
import GeneralInfo from '@/components/GeneralInfo'
import PhysicalTransfer from '@/components/PhysicalTransfer'
import DigitalTransfer from '@/components/DigitalTransfer'
import { mapGetters } from "vuex"

export default {
  name: 'submit',
  components: {
    SubmitterInfo: SubmitterInfo,
    GeneralInfo: GeneralInfo,
    DigitalTransfer: DigitalTransfer,
    PhysicalTransfer: PhysicalTransfer
  },
  computed: {
      ...mapGetters(["hasError", "error","digitalTransfer", "physicalTransfer"]),
      showDigital() {
        return this.digitalTransfer === true
      },
      showPhysical() {
        return this.physicalTransfer === true
      }
  },
  created: function () {
    // see if the auth user cookie is set. If so, this page 
    // is being called as a redirect from NetBadge. All prior state
    // kept in vuex is gone. Need to restore key state from other cookie
    let authUser = this.$cookies.get("archives_xfer_user")
    if (authUser) {
      let settings = this.$cookies.get("archives_xfer_settings")
      this.$store.commit("setUVA", true)
      this.$store.commit("setPhysicalTransfer", settings.physical)
      this.$store.commit("setDigitalTransfer", settings.digital)
      this.$store.commit("setUser",authUser)
      this.$cookies.remove("archives_xfer_user")
      this.$cookies.remove("archives_xfer_settings")
    }
    this.$store.dispatch('getGenres')
    this.$store.dispatch('getUploadID')
  }
}
</script>

<style scoped>
a {
   color: cornflowerblue;
   font-weight: 500;
   text-decoration: none;
}
a:hover {
 text-decoration: underline;  
}
div.contact-small {
  font-size: 0.9em;
  color: #666;
  border-bottom: 1px dashed #EB5F0C;
  padding-bottom: 5px;
  margin-bottom: 15px;
}
div.contact-small p {
  margin: 0;
}
div.submit {
  padding: 30px 50px 250px 50px;
  height: 100%;
  background: white;
  min-width: 1000px;
  width: 75%;
  margin: 0 auto;
  border-right: 1px solid #dfdacc;
  border-left: 1px solid #dfdacc;
  font-weight: 400;
  font-size: 0.9em;
  color: #666;
}
</style>
