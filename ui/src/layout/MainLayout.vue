<template>
  <div class="main">
    <div class="header-container">
      <Header></Header>
    </div>
    <div class="content-container">
      <router-view></router-view>
    </div>

    <Update></Update>

    <!-- <div class="footer-container">
      <Footer></Footer>
    </div> -->
  </div>
</template>

<script>
import Header from "./Header";
import Update from "./Update.vue";

export default {
  name: 'MainLayout',
  components: {
    Header, Update,
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

<style lang="less">
  @import './MainLayout.less';
  @import './Zui.less';
</style>
<style lang="less" scoped>
.main {
  display:flex;
  flex-flow: column;
  height: 100%;

  .header-container {
    height: 50px;
    color: #fff;
    background-color: #0056c9;
    border-bottom: 1px solid #e9f2fb;
  }
  .footer-container {
    height: 40px;
    background: #fff;
    border-top: 1px solid #eff1f7;
  }
  .content-container {
    padding: 0 15px 15px;
  }
}
</style>
