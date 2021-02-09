<template>
  <div>
    <a-menu
        :default-selected-keys="['mine']"
        :selected-keys="[selectedKey]"
        :open-keys.sync="openKeys"
        mode="horizontal"
        @click="handleClick"
    >
      <a-menu-item key="mine/list">
        <Icon type="database" :style="{fontSize: '16px'}" />{{$t('msg.mine')}}
      </a-menu-item>
      <a-menu-item key="buildin/config/list">
        <Icon type="build" :style="{fontSize: '16px'}" />{{$t('msg.buildin')}}
      </a-menu-item>
    </a-menu>
  </div>
</template>

<script>
import {Icon} from 'ant-design-vue'

export default {
  name: 'Navbar',
  components: {
    Icon,
  },
  data () {
    return {
      current: [],
      openKeys: [],
    };
  },
  computed: {
    selectedKey: function() {
      return this.$route.path.split('/')[2] === 'mine' ? 'mine/list' : 'buildin/config/list';
    }
  },
  watch: {
    openKeys(val) {
      console.log('openKeys', val);
    },
  },
  methods: {
    handleClick (e) {
      console.log('handleClick', e, this.$route.path, e.key)
      if (e.key.indexOf('buildin') > -1) this.openKeys = ['buildin']

      const path = '/data/' + e.key
      if (this.$route.path != path) this.$router.push(path);
    },
  }
}
</script>

<style lang="less" scoped>
</style>
