import Vue from 'vue'
import App from './App.vue'

// import { Button, message } from 'ant-design-vue';

import VueI18n from 'vue-i18n'
import ConfigProvider from "ant-design-vue/lib/config-provider";
import Button from "ant-design-vue/lib/button";
import 'ant-design-vue/lib/button/style';

import Menu from "ant-design-vue/lib/menu";
import SubMenu from "ant-design-vue/lib/menu";
import MenuItem from "ant-design-vue/lib/menu";
import 'ant-design-vue/lib/menu/style';

import Icon from "ant-design-vue/lib/icon";

import zhCN from './assets/lang/zh-CN'
import router from "./router"

Vue.config.productionTip = false

Vue.component(Button)
Vue.use(VueI18n)
Vue.use(ConfigProvider)
Vue.use(Button);
Vue.use(Menu, SubMenu, MenuItem)
Vue.use(Icon)

const i18n = new VueI18n({
  locale: 'zh-CN',
  messages: {
    'zh-CN': { ...zhCN }
  },
});

new Vue({
  router,
  i18n,
  render: h => h(App),
}).$mount('#app')
