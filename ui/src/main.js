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

import FormModel from "ant-design-vue/lib/form-model";
import Input from "ant-design-vue/lib/input";
import Select from "ant-design-vue/lib/select";
import 'ant-design-vue/lib/form/style';
import 'ant-design-vue/lib/form-model/style';
import 'ant-design-vue/lib/input/style';
import 'ant-design-vue/lib/select/style';

import Popconfirm from "ant-design-vue/lib/popconfirm";
import 'ant-design-vue/lib/popconfirm/style';

import Modal from "ant-design-vue/lib/modal";
import 'ant-design-vue/lib/modal/style';

import InputNumber from "ant-design-vue/lib/input-number";
import 'ant-design-vue/lib/input-number/style';

import Switch from "ant-design-vue/lib/switch";
import 'ant-design-vue/lib/switch/style';

import Radio from "ant-design-vue/lib/radio";
import 'ant-design-vue/lib/radio/style';

import Table from "ant-design-vue/lib/table";
import 'ant-design-vue/lib/table/style';

import Tabs from "ant-design-vue/lib/tabs";
import 'ant-design-vue/lib/tabs/style';

import Tree from "ant-design-vue/lib/tree";
import 'ant-design-vue/lib/tree/style';

import Tag from "ant-design-vue/lib/tag";
import Divider from "ant-design-vue/lib/divider";
import Icon from "ant-design-vue/lib/icon";

import Col from "ant-design-vue/lib/col";
import Row from "ant-design-vue/lib/row";

import Spin from "ant-design-vue/lib/spin";
import 'ant-design-vue/lib/spin/style';

import Popover from "ant-design-vue/lib/popover";
import 'ant-design-vue/lib/popover/style';

import zhCN from './assets/lang/zh-CN'
import router from "./router"

Vue.config.productionTip = false

Vue.use(VueI18n)
Vue.use(ConfigProvider)
Vue.use(Modal)
Vue.use(Menu)
Vue.use(SubMenu)
Vue.use(MenuItem)
Vue.use(FormModel)
Vue.use(Input)
Vue.use(Select)
Vue.use(Button)
Vue.use(Popconfirm)
Vue.use(Table)
Vue.use(Tag)
Vue.use(Icon)
Vue.use(Divider)
Vue.use(Tree)
Vue.use(Tabs)
Vue.use(Row)
Vue.use(Col)
Vue.use(InputNumber)
Vue.use(Switch)
Vue.use(Radio)
Vue.use(Spin)
Vue.use(Popover)

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
