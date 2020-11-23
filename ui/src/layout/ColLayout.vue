<template>
  <div class="all">
    <div class="left">
      <Menu></Menu>
      <div class="sync">
        <a-button @click="syncData" size="small" type="primary">同步数据</a-button>
      </div>
    </div>
    <div class="content">
      <router-view />
    </div>
  </div>
</template>

<script>
import Menu from "./Menu";
import {syncData} from "../api/manage";

export default {
  name: 'ColLayout',
  components: {
    Menu,
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
    syncData() {
      console.log("syncData")
      syncData().then(json => {
        console.log('syncData', json)
        if (json.code == 1) {
          this.$notification['success']({
            message: '提示',
            description:
                '成功同步资源和数据！',
            placement: 'bottomLeft',
            duration: 0.5,
          });
        }
      })
    },
  }
}
</script>

<style lang="less" scoped>
.all {
  display:flex;
  height: 100%;
  width: 100%;

  .left {
    width: 200px;
    border-right: 1px solid #e9f2fb;
    position: relative;
  }
  .content {
    flex: 1;
    padding: 0 10px;
  }
}

.sync {
  margin-top: 10px;
  text-align: center;
}
</style>
