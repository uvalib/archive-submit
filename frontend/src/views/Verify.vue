<template>
   <div class="verify content">
      <template v-if="state === 'verify'">
          <h3>Verifying</h3>
         <p>
            Verifying your account...
         </p>
      </template>
      <template v-else-if="state === 'failed'">
         <h3>Failed verification</h3>
         <p>
            Account verification failed with the following response:
            <br/><span class="error-message">{{ verifyError }}</span>
         </p>
      </template>
      <template v-else-if="state === 'verified'">
         <h3>Account Verified</h3>
         <p>You account has been verified. Click the 'Continue' button below to access the transfer form</p>
         <button @click="continueClicked" class="pure-button pure-button-primary">Continue</button>
      </template>
   </div>
</template>

<script>
import axios from "axios"
export default {
   name: "verify",
   components: {
   },
   data: function() {
      return {
         state: "verify",
         verifyError: null
      };
   },
   created: function () {
      let token = this.$route.params.token
      axios
         .post("/api/verify/"+token)
         .then((response) => {
            this.state = "verified"
            this.$store.commit("setUser", response.data)
         })
         .catch(error => {
            this.state = "failed"
            this.verifyError =  error.response.data
         })
   },
   methods: {
      continueClicked() {
         this.$router.push("/submit")
      }
   }
};
</script>

<style scoped>
.error-message {
   width: 100%;
   text-align: center;
   color: firebrick;
   font-style: italic;
   min-height: 25px;
}
div.controls {
   text-align: right;
   width: 100%;
   padding: 15px;
}
</style>
