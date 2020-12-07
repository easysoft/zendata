<template>
  <div>
    <a-menu
        :default-selected-keys="['mine']"
        :open-keys.sync="openKeys"
        mode="inline"
        @click="handleClick"
    >
      <a-menu-item key="mine/list">
        {{ $t('msg.mine') }}
      </a-menu-item>

      <a-sub-menu key="buildin" @titleClick="titleClick">
        <span slot="title"><span>{{ $t('msg.buildin') }}</span></span>
          <a-menu-item key="buildin/config/list">
            {{ $t('msg.config') }}
          </a-menu-item>
          <a-menu-item key="buildin/ranges/list">
            {{ $t('msg.ranges') }}
          </a-menu-item>
          <a-menu-item key="buildin/instances/list">
            {{ $t('msg.instances') }}
          </a-menu-item>

          <a-menu-item key="buildin/text/list">
            {{ $t('msg.text') }}
          </a-menu-item>
          <a-menu-item key="buildin/excel/list">
            {{ $t('msg.excel') }}
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
      openKeys: [],
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
      if (e.key.indexOf('buildin') > -1) this.openKeys = ['buildin']

      const path = '/data/' + e.key
      if (this.$route.path != path) this.$router.push(path);
    },
    titleClick (e) {
      console.log('titleClick', e)
    },
  }
}
</script>

<style lang="less" scoped>

</style>
