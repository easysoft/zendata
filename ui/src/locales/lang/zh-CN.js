import antd from 'ant-design-vue/es/locale-provider/zh_CN'
import momentCN from 'moment/locale/zh-cn'

const components = {
  antLocale: antd,
  momentName: 'zh-cn',
  momentLocale: momentCN
}

const locale = {
  'site.title': 'ZenData数据生成工具',

  'menu.data.list': '数据列表',
  'menu.data.edit': '数据编辑',
  'menu.config.list': '字段列表',
  'menu.config.edit': '字段编辑',
  'menu.ranges.list': '序列列表',
  'menu.ranges.edit': '序列编辑',
  'menu.instances.list': '实例列表',
  'menu.instances.edit': '实例编辑',
  'menu.excel.list': '表格列表',
  'menu.excel.edit': '表格编辑',
  'menu.text.list': '文本列表',
  'menu.text.edit': '文本编辑',
}

export default {
  ...components,
  ...locale
}
