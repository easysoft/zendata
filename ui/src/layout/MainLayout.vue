<template>
  <div>
    <div class="header">HEADER</div>
    <router-view></router-view>
    <div>FOOTER</div>
  </div>
</template>

<script>
export default {
  name: 'MainLayout',
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
  @import './MainLayout.less';
</style>
