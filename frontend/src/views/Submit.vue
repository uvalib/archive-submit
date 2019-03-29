<template>
  <div class="submit">
    <h2>Archives Records Transfer Form</h2>
    <vue-dropzone :useCustomSlot=true id="customdropzone" :options="dropzoneOptions">
      <div class="dropzone-custom">
        <div class="upload title">Drag and drop to upload content</div>
        <div class="upload subtitle">or click to select a file from your computer</div>
      </div>
    </vue-dropzone>
  </div>
</template>

<script>
import vue2Dropzone from 'vue2-dropzone'
import 'vue2-dropzone/dist/vue2Dropzone.min.css'
export default {
  name: 'submit',
  components: {
    vueDropzone: vue2Dropzone
  },
  data: function () {
    return {
      dropzoneOptions: {
        url: '/api/submit',
        createImageThumbnails: false,
        // timeout: null,    no tmimeouts
        // addRemoveLinks: true,
        maxFilesize: null,
        chunking: true,
        chunkSize: 10000000, // bytes = 10Mb,
        previewTemplate: this.template()
      }
    }
  },
  methods: {
    template: function () {
        return `<div class="dz-preview dz-file-preview" style="width:100px; margin:  5px;">
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
            </div>
        `;
      },
  }
}
</script>

<style scoped>

div.submit {
  padding: 15px 25px;
  min-height: 600px;
}
h2 {
  margin: 0;
  color:#51822F;
  font-weight: 500;
}
div.dropzone-custom {
  color: #666;
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
