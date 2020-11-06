<template>
  <div>
    <a-menu
        :default-selected-keys="['item1']"
        :open-keys.sync="openKeys"
        mode="inline"
        @click="handleClick"
    >
      <a-menu-item key="item1">
        我的数据
      </a-menu-item>

      <a-sub-menu key="sub2" @titleClick="titleClick">
        <span slot="title"><span>內置数据</span></span>
          <a-menu-item key="sub2-1">
            Excel数据
          </a-menu-item>
          <a-menu-item key="sub2-2">
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
      openKeys: ['sub2'],
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
    handleClick () {
      console.log('handleClick')
    },
    titleClick () {
      console.log('titleClick')
    },
  }
}
</script>

<style lang="less" scoped>

</style>
