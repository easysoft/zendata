<template>
  <div>
    <a-menu
        :default-selected-keys="['mine']"
        :open-keys.sync="openKeys"
        mode="inline"
        @click="handleClick"
    >
      <a-menu-item key="mine">
        我的数据
      </a-menu-item>

      <a-sub-menu key="buildin" @titleClick="titleClick">
        <span slot="title"><span>內置数据</span></span>
          <a-menu-item key="buildin/excel">
            Excel数据
          </a-menu-item>
          <a-menu-item key="buildin/yaml">
            YAML数据
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
      console.log('handleClick', e, this.$route.path)
      if (e.key == 'mine' && this.$route.path != '/data/mine/index') {
        this.$router.push('/data/mine/index');
      }
    },
    titleClick (e) {
      console.log('titleClick', e)
    },
  }
}
</script>

<style lang="less" scoped>

</style>
