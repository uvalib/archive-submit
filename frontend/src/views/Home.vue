<template>
   <div class="home content">
      <h2>
         <span>University Archives Records Transfer Form</span>
         <span class="admin">
            <router-link to="/admin">
               <button class="pure-button pure-button-primary">Admin Access</button>
            </router-link>
         </span>
      </h2>
      <div class="contact">
         <p>University Archives</p>
         <p>Albert & Shirley Small Special Collections Library</p>
         <p>P.O. Box 400110</p>
         <p>Charlottesville, VA, 22904-4110</p>
         <p>
            <span class="label">Contact:</span>Bethany Anderson, University Archivist
         </p>
         <p>
            <span class="label">Phone:</span>(434) 982-2980
         </p>
         <p>
            <span class="label">Email:</span>
            <a href="mailto:bga3d@virginia.edu">bga3d@virginia.edu</a>
         </p>
      </div>
      <div class="info">
         <p>
            Donors or their representatives should complete this form when transferring paper-based
            and/or digital administrative records to the University Archives. When possible,
            send this form (and a box/folder title list for records over 1 cubic foot) electronically
            prior to sending records. We will then advise you on the appropriate transfer method.
         </p>
         <p>
            Digital records may be transferred through this form. Please zip files prior to uploading
            for transfer. Please do not transfer files via this form that contain highly sensitive or
            personally identifiable information (PII), such as social security numbers, credit card numbers,
            or information covered by right to privacy laws such as FERPA or HIPPA. For additional information,
            see
         </p>
         <ul>
            <li>
               <a
                  href="https://uvapolicy.virginia.edu/policy/IRM-012"
                  target="_blank"
               >IRM-012: Privacy and Confidentiality of University Information</a>
            </li>
            <li>
               <a
                  href="https://uvapolicy.virginia.edu/policy/IRM-003"
                  target="_blank"
               >IRM-003: Data Protection of University Information</a>
            </li>
         </ul>
         <p>Contact the University Archivist in advance to discuss records containing PII prior to arranging transfer.</p>
         <p>
            <b>
               It is necessary to consult with the University Archivist in advance of any records transfers to
               ensure that the material is appropriate and in a condition for archival retention.
            </b>
            If you have any questions about this form, please contact the University Archivist.
         </p>
      </div>
      <div class="access-type">
         <h3>Access the transfer form</h3>
         <label>
            <input @click="uvaStatusClick(true)" type="radio" name="is_uva" id="uva" value="no">
            I am UVA faculty, staff, or student. (You will be asked to verify your identity using
            <a
               href="http://itc.virginia.edu/netbadge/"
            >NetBadge</a>.)
         </label>
         <label>
            <input
               @click="uvaStatusClick(false)"
               type="radio"
               name="is_uva"
               id="guest"
               value="yes"
               checked
            >
            I am not affiliated with UVA.
         </label>
         <button
            @click="continueClicked"
            id="continue-btn"
            class="pure-button pure-button-primary"
         >Continue</button>
      </div>
   </div>
</template>

<script>
export default {
   name: "home",
   methods: {
      uvaStatusClick: function(status) {
         this.$store.commit("setUVA", status);
      },
      continueClicked: function(/*event*/) {
         this.$store.commit("clearUser")
         if (this.$store.getters.isUVA == false) {
            this.$router.push("verify");
         } else {
            window.location.href = "/authenticate"
         }
      }
   }
};
</script>

<style scoped>
div.contact {
   color: #444;
   padding: 10px 0 15px 0;
   border-bottom: 1px dashed #eb5f0c;
   margin-bottom: 25px;
}
div.contact p {
   margin: 0;
   padding: 0;
}
td.label {
   font-weight: bold;
   padding-right: 10px;
   text-align: right;
}
#continue-btn {
   float: right;
   position: relative;
   top: -28px;
}
span.admin {
   font-size: 12px;
   position: absolute;
   right: 0;
}

div.home {
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
   color: #333;
}
div.info {
   margin-bottom: 15px;
}
span.label {
   font-weight: bold;
   color: #666;
   width: 75px;
   margin-right: 10px;
   text-align: right;
   display: inline-block;
}
h3 {
   color: #eb5f0c;
   border-bottom: 1px dashed #eb5f0c;
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
a {
   color: cornflowerblue;
   font-weight: 500;
   text-decoration: none;
}
a:hover {
   text-decoration: underline;
}
</style>