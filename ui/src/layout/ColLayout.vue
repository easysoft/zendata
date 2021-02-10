<template>
  <div class="content">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'ColLayout',
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
    padding: 0 10px;
  }
}

.sync {
  margin-top: 10px;
  text-align: center;
}
</style>
