<template>
   <AccordionContent title="Digital Records Transfer">
      <div class="data-form">
         <div class="pure-u-1-1 gap">
            <label>
               Describe Technical Information
               <span class="note">(e.g., software that created files, OS, hardware, naming conventions, and original location).</span>
            </label>
            <textarea class="pure-u-1-1" id="dig-tech-description" v-model="description"></textarea>
         </div>
         <div class="pure-u-1-1">
            <label for="dig-date-range">Date Range of Files</label>
            <input id="dig-date-range" class="pure-u-1-1" type="text" v-model="dateRange">
         </div>
         <div class="pure-u-1-1 gap">
            <label for="dig-types">Record Types <span class="note">(check all that apply)</span></label>
            <div class="choices">
               <span v-for="rt in digitalRecordTypes" :key="rt.id">
                  <label class="pure-checkbox inline">
                     <input type="checkbox" name="dig-record-type" :value="rt.id" v-model="selectedTypes">
                     {{ rt.name }}<span class="note">{{rt.description}}</span>
                  </label>
               </span>
            </div>
         </div>
         <vue-dropzone :useCustomSlot=true id="customdropzone" 
               :options="dropzoneOptions" 
               v-on:vdropzone-sending="sendingEvent"
               v-on:vdropzone-success="fileAddedEvent"
               v-on:vdropzone-removed-file="fileRemovedEvent">
            <div class="dropzone-custom">
               <div class="upload title">Drag and drop to upload content</div>
               <div class="upload subtitle">or click to select a file from your computer</div>
            </div>
         </vue-dropzone>
         <div class="total-size">
            <span><b>Total Upload Size: </b></span><span>{{uploadSize}}</span>
         </div>
      </div>
   </AccordionContent>
</template>

<script>
import AccordionContent from '@/components/AccordionContent'
import vue2Dropzone from 'vue2-dropzone'
import 'vue2-dropzone/dist/vue2Dropzone.min.css'
import { mapFields } from 'vuex-map-fields'

export default {
   components: {
      AccordionContent: AccordionContent,
      vueDropzone: vue2Dropzone,
   },
   data: function () {
      return {
         destroyStarted: false,
         dropzoneOptions: {
            url: '/api/upload',
            createImageThumbnails: true,
            maxFilesize: null,
            chunking: true,
            chunkSize: 10000000, // bytes = 10Mb,
             addRemoveLinks: true 
         }
      }
   },
   beforeDestroy() {
      this.destroyStarted = true
   },
   computed: {
      digitalRecordTypes() {
         return this.$store.state.digitalRecordTypes
      },
      uploadSize() {
         return this.$store.getters.digitalUploadSize
      },
      ...mapFields([
         'digital.uploadedFiles',
         'digital.summary',
         'digital.description',
         'digital.dateRange',
         'digital.selectedTypes',
      ])
   },
   methods: {
      fileAddedEvent (file) {
         this.$store.commit("addUploadedFile",file)
      },
      fileRemovedEvent (file) {
         if (this.destroyStarted === false ) {
            this.$store.dispatch("removeUploadedFile",file)
         }
      },
      sendingEvent (file, xhr, formData) {
         formData.append('identifier', this.$store.getters.submissionID);
      },
  }
}
</script>

<style scoped>
div.gap {
   margin: 10px 0;
}
div.dropzone-custom {
  color: #666;
}
.inline {
   display: inline-block;
   margin-left:15px;
}
.inline .note {
   margin-left: 5px;
}
.note {
   color: #999;
   font-size: 0.85em;
}
div.upload.title {
  font-weight: bold;
  font-size: 1.5em;
  margin-top: 15px;
}
div.upload.subtitle {
  font-weight: 500;
  margin-top: 5px;
}
</style>