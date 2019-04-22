<template>
  <div class="submit content">
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
        <template v-if="digitalTransfer">
          <DigitalTransfer/>
        </template>
        <template v-if="physicalTransfer">
          <PhysicalTransfer/>
        </template>
      </fieldset>
    </form>
    <div class="error">{{error}}</div>
    <span @click="submitClicked" class="submit pure-button pure-button-primary">Submit</span>
  </div>
</template>

<script>
import SubmitterInfo from '@/components/SubmitterInfo'
import GeneralInfo from '@/components/GeneralInfo'
import PhysicalTransfer from '@/components/PhysicalTransfer'
import DigitalTransfer from '@/components/DigitalTransfer'
import { mapState } from "vuex"
import axios from 'axios'

export default {
  name: 'submit',
  components: {
    SubmitterInfo: SubmitterInfo,
    GeneralInfo: GeneralInfo,
    DigitalTransfer: DigitalTransfer,
    PhysicalTransfer: PhysicalTransfer
  },
  computed: {
       ...mapState({
         error: state => state.error,
         user: state => state.user,
         accession: state => state.accession,
         digital: state => state.digital,
         physical: state => state.physical,
         digitalTransfer: state => state.digitalTransfer,
         physicalTransfer: state => state.physicalTransfer,
      }),
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
    } else {
      if (this.user.id === 0 ) {
        this.$router.push("/")
      }
    }
    this.$store.dispatch('getGenres')
    this.$store.dispatch('getRecordTypes')
    this.$store.dispatch('getSubmissionID')
  },
  methods: {
    submitClicked() {
      // clean up data from store (put into heirarchy / remove some fields) and send to server as an accession
      let json = this.accession
      json.user = this.user
      json.digitalTransfer = this.digitalTransfer
      json.physicalTransfer = this.physicalTransfer
      if (this.digitalTransfer) {
        json.digital = this.digital
      }
      if (this.physicalTransfer) { 
        // NOTE: Clone the physical data here because we will convert some 
        // of it to int or bool to help back end processing. Don't want to mess 
        // up the store data
        json.physical = {}
        Object.assign(json.physical, this.physical)
        json.physical.transferMethod = parseInt(json.physical.transferMethod, 10)
        json.physical.hasDigital = (json.physical.hasDigital === "1")
        json.physical.hasSoftware = (json.physical.hasSoftware === "1")
        if ( !json.physical.hasDigital) {
          delete json.physical.techInfo 
          delete json.physical.mediaCarriers 
          delete json.physical.mediaCount 
          delete json.physical.hasSoftware 
        }
      }
      // console.log(json)
      axios.post("/api/submit", json).then((/*response*/)  =>  {
        this.$store.commit("clearSubmissionData") 
        this.$router.push("thanks")
      }).catch((error) => {
        this.$store.commit("setError",error.response.data) 
      })
    }
  }
}
</script>

<style scoped>
div.error {
  font-style: italic;
  color: firebrick;
  padding: 5px 00 15px;
}
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

</style>
