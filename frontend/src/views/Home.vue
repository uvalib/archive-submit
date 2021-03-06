<template>
   <div class="home content">
      <span class="admin">
         <button @click="adminClicked" class="pure-button pure-button-primary">Admin Access</button>
      </span>
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
            and/or digital administrative records to the University Archives. 
            This form is <b>only</b> for the transfer of university records, not personal archives 
            donations; to donate personal archives, contact the University Archivist. If you have 
            specific questions about records retention, please contact 
            <a href="https://recordsmanagement.virginia.edu/" target="_blank">Records Management</a>.
            When possible, send this form (and a box/folder title list for records over 1 cubic foot) 
            electronically prior to sending records. We will then advise you on the appropriate 
            transfer method.
         </p>
         <p>
            Digital records may be transferred through this form. 
            <span  id="popover-trigger-disabled">
               Please compress files
               <div id="popover" class="hidden">
                  <span>Compression Info</span>
                  <i @click="togglePopover" id="close-popover" class="fas fa-times-circle"></i>
                  <div class="popover-content">
                     Files compressed with zip, gzip or tar are accepted. 
                     For more details, click the link for your system.
                  </div>
                  <a class="popover-link" 
                     href="https://support.microsoft.com/en-us/help/14200/windows-compress-uncompress-zip-files" 
                     target="_bank">Windows</a>
                  <a class="popover-link" 
                     href="http://osxdaily.com/2012/01/10/how-to-zip-files-in-mac-os-x/" 
                     target="_bank">Mac</a>
                  <a class="popover-link" 
                     href="https://www.howtogeek.com/248780/how-to-compress-and-extract-files-using-the-tar-command-on-linux/" 
                     target="_bank">Linux Command-Line</a>
               </div>
            </span> 
            prior to uploading
            for transfer. Please do not transfer files via this form that contain highly sensitive data or
            personally identifiable information (PII), such as social security numbers, credit card numbers,
            or information covered by right to privacy laws such as FERPA or HIPPA. For additional information,
            see
         </p>
         <ul>
            <li>
               <a href="https://uvapolicy.virginia.edu/policy/IRM-012" target="_blank">
                  IRM-012: Privacy and Confidentiality of University Information
               </a>
            </li>
            <li>
               <a href="https://uvapolicy.virginia.edu/policy/IRM-003" target="_blank">
                  IRM-003: Data Protection of University Information
               </a>
            </li>
         </ul>
         <p>Contact the University Archivist in advance to discuss records containing PII 
            or highly sensitive data prior to arranging transfer.</p>
         <p>
            <b>
               It is necessary to consult with the University Archivist in advance of any records transfers to
               ensure that the material is appropriate and in a condition for archival retention.
               ** Please note that all records need to be assessed for long-term archival value and 
               transferring records via this form does not automatically guarantee that the content 
               will be preserved.
            </b>
            If you have any questions about this form or transferring records, please contact the University Archivist.
         </p>
      </div>
      <div class="access-type">
         <h3>Access the transfer form</h3>
         <h4>Affiliation</h4>
         <label>
            <input @click="uvaStatusClick(true)" type="radio" name="is_uva" value="no">
            I am UVA faculty, staff, or student. (You will be asked to verify your identity using
            <a href="http://itc.virginia.edu/netbadge/">NetBadge</a>.)
         </label>
         <label>
            <input @click="uvaStatusClick(false)" type="radio" name="is_uva"  value="yes" checked>
            I am not affiliated with UVA.
         </label>
         <h4 class="gap">Type of records transfer (check one or both)</h4>
         <label>
            <input type="checkbox" id="physical">
            I want to transfer physical materials (paper-based, other analog, and/or digital media carriers)
         </label>
         <label>
            <input type="checkbox" id="digital" checked>
            I want to transfer born-digital materials via upload through this form
         </label>
         <div class="error-message">{{ error }}</div>
         <button id="continue-btn" @click="continueClicked" class="pure-button pure-button-primary">Continue</button>
      </div>
   </div>
</template>

<script>
import { mapState } from "vuex"
export default {
   name: "home",
   computed: {
      ...mapState({
         error: state => state.error,
      })
   },
   methods: {
      togglePopover: function( event ) {
         event.stopPropagation()
         let po = document.getElementById("popover")
         if (po.classList.contains("hidden")) {
            po.classList.remove("hidden")
         } else  {
            po.classList.add("hidden")
         }
      },
      adminClicked: function() {
         // redirect to the authenticate endpoint (not vue) which is
         // behind NetBadge. If successful, an API token will be generated and
         // stored in an http-only, secure cookie. User will 
         // be redirected to the admin page
         window.location.href = "/authenticate?url=/admin"
      },
      uvaStatusClick: function(status) {
         this.$store.commit("setUVA", status)
      },
      continueClicked: function(/*event*/) {
         let digital = document.getElementById("digital").checked  
         let physical = document.getElementById("physical").checked  
         if (!physical && !digital) {
            this.$store.commit("setError", "Please select at least one type of record to transfer")
            return
         }
         this.$store.commit("transfer/setPhysicalTransfer", physical)
         this.$store.commit("transfer/setDigitalTransfer", digital)
         this.$store.commit("clearUser")
         if (this.$store.state.isUVA == false) {
            this.$router.push("access")
         } else {
            // This redirect to a page not under vue control resets 
            // state held in vuex. Need to persist key bits in a cookie 
            // for retrieval after the auth redirects happen
            this.$cookies.set("archives_xfer_settings", {physical: physical, digital: digital})
            window.location.href = "/authenticate?url=/submit"
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
   right: -1px;
   top: -1px;
}

span.admin button.pure-button.pure-button-primary {
   border: none;
   border-radius: 0 0 0 15px;
   padding: 4px 15px 6px 15px;
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
h4 {
   margin:5px 0;
   color:#555;
}
h4.gap {
   margin-top: 15px;
}
.access-type input {
   margin-right: 10px;
}
.access-type label {
   display: block;
   margin: 0 0 10px 25px;
}
a, #popover-trigger {
   color: cornflowerblue;
   font-weight: 500;
   text-decoration: none;
}
a:hover, #popover-trigger:hover {
   text-decoration: underline;
   cursor: pointer;
}
#popover-trigger  {
   position: relative;
}
.hidden {
   display:none;
}
#close-popover {
   float: right;
   cursor: pointer;
}
a.popover-link {
   color: white;
   margin-left: 20px;
   font-weight: bold;
   display: block;
}
#popover {
   top:25px;
   left:0;
   cursor: default;
   position: absolute;
   background: cornflowerblue;
   font-size: 0.8em;
   color: white;
   padding: 10px;
   border-radius: 5px;
   width: 200px;
   border: 2px solid darkblue;
}
.popover-content {
   border-top: 1px solid white;
   padding-top: 5px;
   margin-top: 2px;
   margin-bottom: 5px;
}
.error-message {
   width: 100%;
   text-align: center;
   color: firebrick;
   font-style: italic;
   min-height: 25px;
}
</style>