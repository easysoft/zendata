<template>
  <div class="header">
    <h2 class="left">
      <a href="https://www.zendata.cn" target="_blank" :title="$t('site.title')"><img src="logo.png" :alt="$t('site.title')"></a>
    </h2>
    <div class="center">
      <Navbar />
    </div>
    <div class="right">
      <div class="dir">
        <div>
          <span><Icon type="folder-open" /> {{$t('msg.workdir')}} |</span>
          <a-button @click="syncData" size="small" type="link" class="btn"><Icon type="sync" :title="$t('msg.help')" />{{ $t('action.import.from.file') }}</a-button>
        </div>
        <code>{{workDir}}</code>
      </div>
      <select-lang :prefixCls="'select-lang'" />
      <a href="https://www.zendata.cn/book/zendata/" target="_blank"><Icon type="question-circle" :title="$t('msg.help')" :style="{fontSize: '18px'}" /></a>
    </div>
  </div>
</template>

<script>
import {Icon} from 'ant-design-vue'
import {getWorkDir, syncData} from "../api/manage";
import {config} from "../utils/vari";
import SelectLang from '../components/SelectLang'
import Navbar from './Navbar';

export default {
  name: 'Header',
  components: {
    SelectLang,
    Icon,
    Navbar,
  },
  data () {
    return {
      workDir: '',
    }
  },
  created () {
    getWorkDir().then(json => {
      console.log('getWorkDir', json)
      const that = this
      that.defs = json.data
      this.workDir = json.workDir
      config.workDir = this.workDir
    })
  },
  methods: {
    syncData() {
      syncData().then(json => {
        if (json.code == 1) {
          this.$notification['success']({
            message: this.$i18n.t('tips.success.to.import'),
            duration: 3,
          });
        }
      })
    },
  }
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
    > a {
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
      div > span {
        opacity: .6;
      }
      .btn {
        font-size: 12px;
        color: #fff;
      }
    }
  }
  .select-lang {
    padding-right: 15px;
    cursor: pointer;
  }
}
</style>
