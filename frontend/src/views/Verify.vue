<template>
   <div class="register">
      <h2>University Archives Records Transfer Form</h2>
      <template v-if="state === 'verify'">
         <p>
            To access the transfer form, please enter your email below.
            <br>If the address has not
            yet been verified, you will be walked through the verification process.
         </p>
         <div class="email-entry pure-form">
            <label>Email Address</label>
            <input type="email" id="email" placeholder="Email">
            <button
               @click="checkEmail"
               type="submit"
               class="pure-button pure-button-primary"
            >Continue</button>
         </div>
      </template>
      <template v-else-if="state === 'verified'">
         <h2>Verified</h2>
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
   </div>
</template>

<script>
import axios from "axios";
import SubmitterInfo from "@/components/SubmitterInfo";
import { mapGetters } from "vuex";
export default {
   name: "register",
   components: {
      SubmitterInfo: SubmitterInfo
   },
   computed: {
      ...mapGetters(["user", "hasError", "error"])
   },
   data: function() {
      return {
         state: "verify"
      };
   },
   methods: {
      verifyClicked() {
         let user = {};
         user.lastName = document.getElementById("lname").value;
         user.firstName = document.getElementById("fname").value;
         user.title = document.getElementById("title").value;
         user.affiliation = document.getElementById("affiliation").value;
         user.email = document.getElementById("email").value;
         user.phone = document.getElementById("phone").value;
         axios
            .post("/api/users", user)
            .then(response => {
               this.$store.commit("setUser", response.data);
               this.state = "submitted";
            })
            .catch(error => {
               this.$store.commit("setError", error.response.data);
            });
      },
      checkEmail() {
         let emailInput = document.getElementById("email");
         let email = emailInput.value;
         axios
            .get("/api/users/lookup?email=" + email)
            .then(response => {
               // TODO handle a case where a user exists, but is not verified
               this.$store.commit("setUser", response.data);
               this.state = "verified";
            })
            .catch((/*error*/) => {
               this.state = "register";
               this.$store.commit("setUserEmail", email);
            });
      }
   }
};
</script>

<style scoped>
div.register {
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
