import Vue from 'vue'
import App from './App.vue'

import Button from "ant-design-vue/lib/button";
import "ant-design-vue/dist/antd.css";

Vue.component(Button.name, Button);
import router from "./router"

Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
