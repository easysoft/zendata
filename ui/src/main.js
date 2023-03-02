import Vue from 'vue'
import App from './App.vue'
import i18n from './locales'
import store from './store/'
import bootstrap from './config/bootstrap'

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

import Tooltip from "ant-design-vue/lib/tooltip";
import 'ant-design-vue/lib/tooltip/style';

import Pagination from "ant-design-vue/lib/pagination";
import 'ant-design-vue/lib/pagination/style';

import Card from "ant-design-vue/lib/card";
import "ant-design-vue/lib/card/style";

import Upload from "ant-design-vue/lib/upload";
import "ant-design-vue/lib/upload/style";

import Drawer from "ant-design-vue/lib/drawer";
import "ant-design-vue/lib/drawer/style";

import Progress from "ant-design-vue/lib/progress";
import "ant-design-vue/lib/progress/style";

import notification from "ant-design-vue/lib/notification";
import 'ant-design-vue/lib/notification/style';
import message from "ant-design-vue/lib/message";
import 'ant-design-vue/lib/message/style';

import VueClipboard from 'vue-clipboard2'

Vue.prototype.$message = message;
Vue.prototype.$notification = notification;

import router from "./router"


Vue.config.productionTip = false

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
Vue.use(Tooltip)
Vue.use(Pagination)
Vue.use(Card)
Vue.use(Upload)
Vue.use(Drawer)
Vue.use(Progress)
Vue.use(VueClipboard)

new Vue({
  router,
  store,
  i18n,
  created: bootstrap,
  render: h => h(App),
}).$mount('#app')
