<template>
  <div class="all">
    <div class="left">
      <Menu></Menu>
    </div>
    <div class="content">
      <router-view />
    </div>
  </div>
</template>

<script>
import Menu from "./Menu";

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
  }
  .content {
    flex: 1;
    padding: 0 10px;
  }
}
</style>
