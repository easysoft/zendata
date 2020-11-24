<template>
  <div>
    <a-menu
        :default-selected-keys="['mine']"
        :open-keys.sync="openKeys"
        mode="inline"
        @click="handleClick"
    >
      <a-menu-item key="mine/list">
        我的数据
      </a-menu-item>

      <a-sub-menu key="buildin" @titleClick="titleClick">
        <span slot="title"><span>內置数据</span></span>
          <a-menu-item key="buildin/config/list">
            字段
          </a-menu-item>
          <a-menu-item key="buildin/ranges/list">
            序列
          </a-menu-item>
          <a-menu-item key="buildin/instances/list">
            实例
          </a-menu-item>

          <a-menu-item key="buildin/text/list">
            文本
          </a-menu-item>
          <a-menu-item key="buildin/excel/list">
            表格
          </a-menu-item>
      </a-sub-menu>
    </a-menu>
  </div>
</template>

<script>
export default {
  name: 'Menu',
  components: {
  },
  data () {
    return {
      current: [],
      openKeys: ['buildin'],
    };
  },
  watch: {
    openKeys(val) {
      console.log('openKeys', val);
    },
  },
  computed: {
  },
  created () {
  },
  mounted () {
  },
  methods: {
    handleClick (e) {
      console.log('handleClick', e, this.$route.path, e.key)
      // if (e.key == 'mine' && this.$route.path != '/data/mine/index') {
      //   this.$router.push('/data/mine/index');
      // } else {
        const path = '/data/' + e.key
        if (this.$route.path != path) this.$router.push(path);
      // }
    },
    titleClick (e) {
      console.log('titleClick', e)
    },
  }
}
</script>

<style lang="less" scoped>

</style>
