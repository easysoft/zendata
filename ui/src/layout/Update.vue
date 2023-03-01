<template>
  <div class="right-top-update-main">
    <div v-if="version" class="version">V{{version}}</div>

    <a-modal :title="$t('update.title')"
             :visible="isVisible"
             @cancel="onCancel"
             :maskClosable="false"
             class="update-modal">
      <div>
        {{ $t('update.new.pre') }}<b>{{newVersion}}</b>{{ $t('update.new.suf') }}
      </div>

      <div v-if="downloadingPercent > 0">
        <a-progress :percent="downloadingPercent" />
      </div>

      <div v-if="errMsg" class="errors">
        <div class="border">{{ $t('update.failed') }}</div>
        <div>{{errMsg}}</div>
      </div>

      <template #footer>
        <div v-if="!updateSuccess">
          <a-button @click="update" type="primary">{{ $t('update.update') }}</a-button>
          <a-button @click="defer">{{ $t('update.notice.tomorrow') }}</a-button>
          <a-button @click="skip">{{ $t('update.notice.skip') }}</a-button>
        </div>

        <div v-if="updateSuccess">
          <a-button @click="rebootNow" type="primary">{{ $t('update.reboot') }}</a-button>
          <a-button @click="rebootLater">{{ $t('update.pending') }}</a-button>
        </div>
      </template>
    </a-modal>
  </div>
</template>

<script>
import {
  electronMsgUpdate,
  electronMsgDownloading,
  skippedVersion, ignoreUtil, electronMsgDownloadSuccess, electronMsgUpdateFail, electronMsgReboot,
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
    let version = null
    let updateSuccess = false
    let errMsg = ''

    return {
      isVisible,
      currVersion,
      newVersion,
      forceUpdate,
      downloadingPercent,
      ipcRenderer,
      isElectron,
      version,
      updateSuccess,
      errMsg,
    }
  },

  created() {
    console.log('created')
    this.isElectron = !!window.require

    if (this.isElectron) {
      const remote = window.require('@electron/remote')
      this.version = remote.getGlobal('sharedObj').version
    }

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
      this.ipcRenderer.on(electronMsgDownloadSuccess, async (event, data) => {
        console.log('md5 checking success msg from electron', data);
        this.updateSuccess = true
      })
      this.ipcRenderer.on(electronMsgUpdateFail, async (event, data) => {
        console.log('downloading fail msg from electron', data);
        this.errMsg = data.err
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
    rebootNow() {
      console.log('rebootNow')
      this.ipcRenderer.send(electronMsgReboot, {})
    },
    rebootLater() {
      console.log('rebootLater')
      this.onCancel()
    },
    defer() {
      console.log('defer')
      setCache(ignoreUtil, Date.now() + 24 * 3600);
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
.right-top-update-main {
  position: absolute;
  right: 5px;
  bottom: 0;
}

.update-modal{
  .ant-modal-footer {
    text-align: center;
  }

  .errors {
    margin-top: 12px;
    .border {
      margin-bottom: 3px;
      font-weight: bolder;
    }
  }
}
</style>