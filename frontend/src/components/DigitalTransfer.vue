<template>
   <AccordionContent title="Digital Records Transfer">
      <div class="data-form">
         <div class="pure-u-1-1 gap">
            <label>
               Describe Technical Information
               <span class="note">(e.g., software that created files, OS, hardware, naming conventions, and original location).</span>
            </label>
            <textarea class="pure-u-1-1" id="dig-tech-description"></textarea>
         </div>
         <div class="pure-u-1-2">
            <label for="dig-date-range">Date Range of Files</label>
            <input id="dig-date-range" class="pure-u-23-24" type="text">
         </div>
         <div class="pure-u-1-2">
            <label for="dig-size">Extent of records <span class="note">(in gigabytes)</span></label>
            <input id="dig-size" class="pure-u-23-24" type="text">
         </div>
         <div class="pure-u-1-1 gap">
            <label for="dig-types">Record Types <span class="note">(check all that apply)</span></label>
            <div class="choices">
               <span v-for="rt in digitalRecordTypes" :key="rt.id">
                  <label class="pure-checkbox inline">
                     <input type="checkbox" name="dig-record-type" :value="rt.id">
                     {{ rt.name }}<span class="note">{{rt.description}}</span>
                  </label>
               </span>
            </div>
         </div>
         <input type="hidden" id="submitted-files" name="submitted-files" :value="uploadedFiles">
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
      </div>
   </AccordionContent>
</template>

<script>
import AccordionContent from '@/components/AccordionContent'
import vue2Dropzone from 'vue2-dropzone'
import 'vue2-dropzone/dist/vue2Dropzone.min.css'
import { mapGetters } from "vuex"

export default {
   components: {
      AccordionContent: AccordionContent,
      vueDropzone: vue2Dropzone,
   },
   data: function () {
      return {
         dropzoneOptions: {
            url: '/api/upload',
            createImageThumbnails: false,
            maxFilesize: null,
            chunking: true,
            chunkSize: 10000000, // bytes = 10Mb,
            previewTemplate: this.template()
         }
      }
   },
   computed: {
      digitalRecordTypes() {
         return this.$store.getters.digitalRecordTypes
      },
      uploadedFiles() {
         return this.$store.getters.uploadedFiles
      }
   },
   methods: {
      fileAddedEvent (file) {
         // just adds filename to store list 
         this.$store.commit("addUploadedFile",file.name)
      },
      fileRemovedEvent (file) {
         // makes an ajax call to the service to remove the file
         this.$store.dispatch("removeUploadedFile",file.name)
      },
      sendingEvent (file, xhr, formData) {
         formData.append('identifier', this.$store.getters.uploadID);
      },
      template: function () {
         return `
         <div class="dz-preview dz-file-preview" style="width:100px; margin:  5px;">
            <style type="text/css">
            .custom-remove { z-index:1000;position:absolute;right:5px;top:5px;opacity:0.8;}
            .custom-remove:hover { opacity:1;cursor:pointer; }
            .dz-progress.custom { margin-top:0 !important; top: auto !important; bottom:10px !important; border-radius:0 !important;background:white !important;}
            .dz-upload.custom { background: darkslateblue !important}
            .truncated {display:block;overflow: hidden;text-overflow: ellipsis;white-space: nowrap;}
            p.failure {text-align:center; padding: 1px 4px 0 4px; background: white; color: #a00;
                        border-left: 4px solid rgba(33,150,243,.8);
                        border-right: 4px solid rgba(33,150,243,.8);padding: 2px 0 0 0;}
            </style>
                  <img class="custom-remove" src="remove.png" data-dz-remove/>
                  <div class="dz-image">
                     <div data-dz-thumbnail-bg></div>
                  </div>
                  <div class="dz-details" style="text-align: center;padding: 30px 10px;">
                     <div class="dz-name" style="font-size:12px"><span class="truncated" data-dz-name></span></div>
                     <div class="dz-size" style="font-size:12px"><span data-dz-size></span></div>
                  </div>
                  <div class="dz-progress custom"><span class="dz-upload custom" data-dz-uploadprogress></span></div>
                  <div class="dz-error-message"><span data-dz-errormessage></span></div>
                  <div class="dz-success-mark"><i class="fa fa-check"></i></div>
                  <div class="dz-error-mark" style="top:60%"><p class="failure">FAILED<p></i></div>
         </div>`
      }
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