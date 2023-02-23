<template>
  <div class="right-top-update-main">
    <div></div>

    <a-modal title="升级提醒"
           :visible="isVisible"
           @cancel="onCancel"
           :maskClosable="false"
           class="update-modal">
      <div>
        发现新的版本<b>{{newVersion}}</b>，请确定是否升级。
      </div>
      <div v-if="downloadingPercent > 0">
        <a-progress :percent="downloadingPercent" />
      </div>

      <template #footer>
        <a-button @click="update" type="primary">立即升级</a-button>
        <a-button @click="defer">明天提醒我</a-button>
        <a-button @click="skip">跳过这个版本</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script>
import {
  electronMsgUpdate,
  electronMsgDownloading,
  skippedVersion, ignoreUtil,
} from "../config/settings.js";
  import {getCache, setCache} from "../utils/localCache.js";

export default {
  name: 'Update',
  components: {
  },
  data() {
    let isVisible = false
    let currVersion = ''
    let newVersion = ''
    let forceUpdate = false
    let downloadingPercent = 0
    let isElectron = false
    let ipcRenderer = undefined

    return {
      isVisible,
      currVersion,
      newVersion,
      forceUpdate,
      downloadingPercent,
      ipcRenderer,
      isElectron,
    }
  },

  created() {
    this.isElectron = !!window.require

    if (this.isElectron && !this.ipcRenderer) {
      this.ipcRenderer = window.require('electron').ipcRenderer

      console.log('ipcRenderer', this.ipcRenderer)
      this.ipcRenderer.on(electronMsgUpdate, async (event, data) => {
        console.log('update msg from electron', data)
        this.currVersion = data.currVersionStr
        this.newVersion = data.newVersionStr
        this.forceUpdate = data.forceUpdate

        const skippedVersionVal = await getCache(skippedVersion);
        const ignoreUtilVal = await getCache(ignoreUtil);
        if (skippedVersionVal === this.newVersion || Date.now() < ignoreUtilVal) return;

        this.isVisible = true
      })

      this.ipcRenderer.on(electronMsgDownloading, async (event, data) => {
        console.log('downloading msg from electron', data);
        this.downloadingPercent = Math.round(data.percent * 100);
      })
    }
  },
  mounted() {
  },

  methods: {
    update() {
      console.log('update')
      this.ipcRenderer.send(electronMsgUpdate, {
        currVersion: this.currVersion,
        newVersion: this.newVersion,
        forceUpdate: this.forceUpdate
      })
    },
    defer() {
      console.log('defer')
      setCache(skippedVersion, Date.now() + 24 * 3600);
      this.isVisible = false
    },
    skip() {
      console.log('skip')
      setCache(skippedVersion, this.newVersion);
      this.isVisible = false
    },

    onCancel() {
      console.log('onCancel')
      this.isVisible = false
    },
  }
}


</script>

<style lang="less">
.update-modal{
  .ant-modal-footer {
    text-align: center;
  }
}
</style>