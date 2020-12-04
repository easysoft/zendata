<template>
  <div class="header">
    <h2 class="left">
      <a href="https://www.zendata.cn" target="_blank">ZenData</a>
    </h2>
    <div class="center"></div>
    <div class="right">
      <span class="dir">工作目录：{{workDir}}</span>
      <a href="https://www.zendata.cn/book/zendata/" target="_blank">帮助</a>
    </div>
  </div>
</template>

<script>

import {getWorkDir} from "../api/manage";
import {config} from "../utils/vari";

export default {
  name: 'Header',
  components: {
  },
  data () {
    return {
      workDir: '',
    }
  },
  computed: {
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
  mounted () {
  },
  methods: {
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
    width: 100px;
    color: #fff;
  }
  .center {
    flex: 1;
  }
  .right {
    margin: 0 15px;
    width: 500px;
    text-align: right;

    .dir {
      display: inline-block;
      padding-right: 20px;
    }
  }
}
</style>
