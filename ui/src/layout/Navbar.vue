<template>
  <div class="navbar">
    <a-menu
        :default-selected-keys="['mine']"
        :selected-keys="[selectedKey]"
        :open-keys.sync="openKeys"
        mode="horizontal"
        @click="handleClick"
    >
      <a-menu-item key="/data/mine/list" class="link">
        <Icon type="database" :style="{fontSize: '16px'}" />{{$t('msg.mine')}}
      </a-menu-item>
      <a-menu-item key="/data/buildin/config/list" class="link">
        <Icon type="build" :style="{fontSize: '16px'}" />{{$t('msg.buildin')}}
      </a-menu-item>
      <a-menu-item key="/mock/index" class="link">
        <Icon type="cloud-server" :style="{fontSize: '16px'}" />{{$t('menu.data.mock')}}
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
      console.log(this.$route.path)
      const arr = this.$route.path.split('/')

      if (arr[2] === 'mine') {
        return '/data/mine/list'
      } else if (arr[2] === 'buildin') {
        return '/data/buildin/config/list'
      } else if (arr[1] === 'mock') {
        return '/mock/index'
      }

      return ''
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

      const path =  e.key
      if (this.$route.path != path) this.$router.push(path);
    },
  }
}
</script>

<style lang="less" scoped>
.navbar {
  link {
    cursor: pointer;
  }
}
</style>
