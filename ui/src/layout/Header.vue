<template>
  <div class="header">
    <h2 class="left">
      <a href="https://www.zendata.cn" target="_blank">zd.im</a>
    </h2>
    <div class="center"></div>
    <div class="right">
      <a href="https://www.zendata.cn/book/zendata/" target="_blank">帮助</a>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Header',
  components: {
  },
  data () {
    return {
    }
  },
  computed: {
  },
  created () {
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
    width: 100px;
    text-align: right;
  }
}
</style>
