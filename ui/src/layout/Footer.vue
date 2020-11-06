<template>
  <div class="footer">
    <span class="copyright">
      © 2019- 2020 <a href="http://www.easycorp.cn" target="_blank">易企天创</a>
    </span>
    <span class="icp-info">
      <a href="http://beian.miit.gov.cn/" class="common" target="_blank">鲁ICP备18054969号-14</a>
    </span>
  </div>
</template>

<script>
export default {
  name: 'Footer',
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
.footer {
  height: 39px;
  line-height: 39px;
  text-align: center;
  font-family: 'WenQuanYi Micro Hei';

  span {
    display: inline-block;
    margin: 0 8px;
  }

  .icp-info {
    a {
      font-size: 12px;
    }
  }
}
</style>
