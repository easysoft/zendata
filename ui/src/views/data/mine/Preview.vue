<template>
  <a-card>
    <div slot="title">
      <a-icon type="profile" />
      <span>{{$t('msg.preview')}}</span>
      <a-input v-model="previewUrl" id="url" />&nbsp;&nbsp;
      <a @click="doCopy" :title="$t('action.design')">
        {{ $t('copy.title') }}
      </a> &nbsp;
      <a @click="loadPreviewData" :title="$t('action.design')">
        {{ $t('preview.title') }}
      </a> &nbsp;
    </div>
    <pre v-if="previewData !== null" v-html="previewData" style="margin: 0"></pre>
    <div v-else style="padding: 10px; text-align: center"><a-icon type="loading" /></div>
  </a-card>
</template>

<script>
import {previewDefData} from "../../../api/manage";
import {serverUrl} from '../../../utils/request'

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
    return {
        previewData: null,
        previewUrl: serverUrl + "/data/generate?format=txt&config=" + this.record.referName.replace(/\\/g, "/"),
    };
  },
  mounted: function() {
    this.loadPreviewData();
  },
  methods: {
    loadPreviewData() {
      console.log(this.record)
      this.previewData = null;
      let params = this.getQuery()

      previewDefData(params).then(data => {
        this.previewData = data
      })
    },
    getQuery() {
    let url = decodeURI(this.previewUrl); // 获取url中"?"符后的字串(包括问号)
    url = url.replace(`${serverUrl}/data/generate`, "")

    let query = {};
    if (url.indexOf("?") != -1) {
        const str = url.substr(1);
        const pairs = str.split("&");
        for(let i = 0; i < pairs.length; i ++) {
             const pair = pairs[i].split("=");
            query[pair[0]] = pair[1];
        }
    }
    return query ;  // 返回对象
  },
  doCopy: function () {
    let that = this;
    this.$copyText(this.previewUrl).then(function () {
        that.$notification['success']({
        message: that.$t('copy.success'),
      });
    }, function (e) {
        console.log(e)
    })
    }
  },
  watch: {
    record: function() {
      this.previewUrl = serverUrl + "/data/generate?format=txt&config=" + this.record.referName.replace(/\\/g, "/");
      this.loadPreviewData();
    }
  }
}
</script>

<style>
#url{
    max-width: 30vw;
    margin-left: 20px;
}
</style>