<template>
   <div class="forbidden">
      <h2>University Archives Records Transfer Form</h2>
      <div>
         <h3>Thank you!</h3>
         <p>
            Your transfer request has successfully been submitted and is being processed.
            <br>
            A receipt has been emailed to: {{ user.email}}
            <br>Once processing has completed, you will receive a notification email detailing the status of the transfer.
         </p>
      </div>
   </div>
</template>

<script>
export default {
   name: "thanks",
   computed: {
      user() {
         return this.$store.getters.user
      }
   }
};
</script>

<style scoped>
div.forbidden {
   padding: 15px 25px;
   min-height: 600px;
   background: white;
}
</style>