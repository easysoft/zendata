<template>
  <a-card>
    <span slot="title">
      <a-icon type="profile" />
      <span>{{$t('msg.preview')}}</span>
    </span>
    <pre v-if="previewData !== null" v-html="previewData" style="margin: 0"></pre>
    <div v-else style="padding: 10px; text-align: center"><a-icon type="loading" /></div>
  </a-card>
</template>

<script>
import {previewDefData} from "../../../api/manage";

export default {
  name: 'Preview',
  components: {
  },
  props: {
    record: {
      type: Object,
      required: true
    },
  },
  data: function() {
    return {previewData: null};
  },
  mounted: function() {
    this.loadPreviewData();
  },
  methods: {
    loadPreviewData() {
      console.log(this.record)
      this.previewData = null;
      previewDefData(this.record.id).then(json => {
        this.previewData = json.data
      })
    }
  },
  watch: {
    record: function() {
      this.loadPreviewData();
    }
  }
}
</script>
