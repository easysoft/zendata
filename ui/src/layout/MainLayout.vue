<template>
  <div class="main">
    <div class="header-container">
      <Header></Header>
    </div>
    <div class="content-container">
      <router-view></router-view>
    </div>
    <div class="footer-container">
      <Footer></Footer>
    </div>
  </div>
</template>

<script>
import Header from "./Header";
import Footer from "./Footer";

export default {
  name: 'MainLayout',
  components: {
    Header, Footer
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
</style>
<style lang="less" scoped>
.main {
  display:flex;
  flex-flow: column;
  height: 100%;

  .header-container {
    height: 50px;
    color: #fff;
    background-color: #1183fb;
    border-bottom: 1px solid #e9f2fb;
  }
  .content-container {
    flex: 1;
    display:flex;
    padding: 10px;
    background: #fff;
  }
  .footer-container {
    height: 40px;
    background: #fff;
    border-top: 1px solid #eff1f7;
  }
}
</style>
