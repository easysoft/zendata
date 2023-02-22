<template>
  <div class="header">
    <h2 class="left">
      <a href="https://www.zendata.cn" target="_blank" :title="$t('site.title')">
        <img :src="logoPath" :alt="$t('site.title')">
      </a>
    </h2>
    <div class="center">
      <Navbar />
    </div>
    <div class="right">
      <div class="dir">
        <div>
          <span>
            <Icon type="folder-open" /> {{ $t('msg.workdir') }} |
          </span>
          <a-button @click="syncData" size="small" type="link" class="btn">
            <Icon type="sync" :title="$t('msg.help')" />{{ $t('action.import.from.file') }}
          </a-button>
        </div>
        <code>{{ workDir }}</code>
      </div>
      <select-lang :prefixCls="'select-lang'" />
      <a href="https://www.zendata.cn/book/zendata/" target="_blank">
        <Icon type="question-circle" :title="$t('msg.help')" :style="{ fontSize: '18px' }" />
      </a>
      <div v-if="isElectron" id="windowBtn">
        <span v-if="!fullScreenDef" @click="fullScreen" :title="$t('window.fullscreen')" class="window-btn">
          <Icon type="fullscreen" class="window-btn" :style="{ fontSize: '18px' }" />
        </span>
        <span v-if="fullScreenDef" @click="fullScreen" :title="$t('window.exit_fullscreen')" class="window-btn">
          <Icon type="fullscreen-exit" class="window-btn" :style="{ fontSize: '18px' }" />
        </span>
  
        <span :title="$t('window.minimize')" @click="minimize" class="window-btn">
          <Icon type="minus" class="window-btn" :style="{ fontSize: '18px' }" />
        </span>
  
        <span v-if="maximizeDef" :title="$t('window.restore')" @click="maximize" class="window-btn">
          <Icon type="block" :style="{ fontSize: '18px' }" />
        </span>
        <span v-if="!maximizeDef" :title="$t('window.maximize')" @click="maximize" class="window-btn">
          <Icon type="border" class="window-btn" :style="{ fontSize: '18px' }" />
        </span>
  
        <span :title="$t('window.close')" @click="exit" class="window-btn">
          <Icon type="close" class="window-btn" :style="{ fontSize: '18px' }" />
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import { Icon } from 'ant-design-vue'
import { getWorkDir, syncData } from "../api/manage";
import { config } from "../utils/vari";
import SelectLang from '../components/SelectLang'
import Navbar from './Navbar';
import { getPath } from "@/utils/dom";
import { getElectron } from "@/utils/common";

export default {
  name: 'Header',
  components: {
    SelectLang,
    Icon,
    Navbar,
  },
  data() {
    return {
      workDir: '',
      logoPath: '',
      maximizeDef: true,
      fullScreenDef: false,
      isElectron: getElectron(),
    }
  },
  created() {
    this.logoPath = getPath() + 'logo.png'

    getWorkDir().then(json => {
      this.workDir = json.data
      config.workDir = this.workDir
    })
  },
  methods: {
    syncData() {
      syncData().then(json => {
        if (json.code == 0) {
          this.$router.go(0)

          // this.$notification['success']({
          //   message: this.$i18n.t('tips.success.to.import'),
          //   duration: 3,
          // });
        }
      })
    },
    fullScreen() {
      console.log('fullScreen')
      this.fullScreenDef = !this.fullScreenDef

      const { ipcRenderer } = window.require('electron')
      ipcRenderer.send('electronMsg', 'fullScreen')
    },

    minimize() {
      console.log('minimize')

      const { ipcRenderer } = window.require('electron')
      ipcRenderer.send('electronMsg', 'minimize')
    },
    maximize() {
      console.log('maximize')

      const { ipcRenderer } = window.require('electron')
      ipcRenderer.send('electronMsg', this.maximizeDef ? 'unmaximize' : 'maximize')
      this.maximizeDef = !this.maximizeDef
    },

    exit() {
      console.log('exit')
      const { ipcRenderer } = window.require('electron')
      ipcRenderer.send('electronMsg', 'exit')
    }
  },
}
</script>

<style lang="less" scoped>
.header {
  display: flex;
  height: 49px;
  line-height: 49px;

  a {
    color: #fff;
  }

  .left {
    margin: 0 15px;
    flex: none;

    >a {
      display: block;

      img {
        height: 38px;
        display: block;
        margin-top: 5px;
      }
    }
  }

  .center {
    flex: 1;
  }

  .right {
    margin: 0 15px;
    display: flex;
    align-items: center;

    .dir {
      padding: 0 20px 3px 10px;
      font-size: 12px;
      line-height: 13px;

      div>span {
        opacity: .6;
      }

      .btn {
        font-size: 12px;
        color: #fff;
      }
    }

    .window-btn {
      margin-left: 8px;
    }
  }

  .select-lang {
    padding-right: 15px;
    cursor: pointer;
  }
}</style>
