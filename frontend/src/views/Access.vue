<template>
   <div class="access content">
      <template v-if="state === 'verify'">
         <p>
            To access the transfer form, please enter your email below.
            <br>If the address has not
            yet been verified, you will be walked through the verification process.
         </p>
         <div class="email-entry pure-form">
            <label>Email Address</label>
            <input type="email" id="email" placeholder="Email">
            <button  @click="checkEmail" class="pure-button pure-button-primary">Continue</button>
         </div>
      </template>
      <template v-else-if="state === 'unverified'">
         <h3>Email '{{ user.email }}' exists, but has not yet been verified</h3>
         <p>
            You should have received an email entitled 'UVA Archives Transfer Verification' that contains a verification link.
            <br/><b>You must click this link before you can access the transfer form.</b>
         </p>
         <p> 
            Check your account (including spam folders) for this message, and click the link if it is found.
            <br/>If you do not have the email, click the resend button below and it will be sent out again.
         </p>
         <div class="controls">
            <button @click="resendClicked" class="pure-button pure-button-primary resend">Resend Verification Email</button>
         </div>
      </template>
      <template v-else-if="state === 'register'">
         <h3>This email address has not yet been verified</h3>
         <p>To verify it, please fill out the form below and click 'Verify'.</p>

         <div class="user-verify pure-form pure-g">
            <SubmitterInfo/>
         </div>
         <div class="error-message">{{ error }}</div>
         <div class="controls">
            <button @click="verifyClicked" class="pure-button pure-button-primary verify">Verify</button>
         </div>
      </template>
      <template v-else-if="state === 'submitted'">
         <h3>Thank You!</h3>
         <p>
            Your information has been received. A verification email has been sent to the address you provided.
            <br>Click the link it contains to active your account and access the transfer form.
         </p>
      </template>
       <template v-else-if="state === 'resent'">
         <h3>Verification Email Resent</h3>
         <p>
            You should receive the email in a few minutes.
            <br>Click the link it contains to active your account and access the transfer form.
         </p>
      </template>
   </div>
</template>

<script>
import axios from "axios"
import SubmitterInfo from "@/components/SubmitterInfo"
import { mapState } from "vuex"
export default {
   name: "access",
   components: {
      SubmitterInfo: SubmitterInfo
   },
   computed: {
      ...mapState({
         error: state => state.error,
         user: state => state.user,
      })
   },
   data: function() {
      return {
         state: "verify"
      }
   },
   methods: {
      resendClicked() {
         axios
            .post("/api/resend/verification", {token:this.user.token})
            .then((/*response*/) => {
               this.state = "resent"
            })
            .catch(error => {
               this.$store.commit("setError", error.response.data)
            })
      },
      verifyClicked() {
         let user = {}
         user.lastName = document.getElementById("lname").value
         user.firstName = document.getElementById("fname").value
         user.title = document.getElementById("title").value
         user.affiliation = document.getElementById("affiliation").value
         user.email = document.getElementById("email").value
         user.phone = document.getElementById("phone").value
         axios
            .post("/api/users", user)
            .then(response => {
               this.$store.commit("setUser", response.data)
               this.state = "submitted"
            })
            .catch(error => {
               this.$store.commit("setError", error.response.data)
            })
      },
      checkEmail() {
         let emailInput = document.getElementById("email")
         let email = emailInput.value
         axios
            .get("/api/users/lookup?email=" + email)
            .then(response => {
               let user = response.data 
               this.$store.commit("setUser", user)
               if (user.verified == true) {
                  this.$router.push("submit")
               } else {
                  this.state = "unverified"
               }
            })
            .catch((/*error*/) => {
               this.state = "register"
               this.$store.commit("setUserEmail", email)
            });
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
.email-entry label {
   display: block;
}
#email {
   width: 300px;
   margin-right: 5px;
}
</style>
