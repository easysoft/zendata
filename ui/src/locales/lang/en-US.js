import antdEnUS from 'ant-design-vue/es/locale-provider/en_US'
import momentEU from 'moment/locale/eu'

const components = {
  antLocale: antdEnUS,
  momentName: 'eu',
  momentLocale: momentEU
}

const locale = {
  'site.title': 'ZenData',

  'menu.data.list': 'Data List',
  'menu.data.edit': 'Data Edit',
  'menu.config.list': 'Config List',
  'menu.config.edit': 'Config Edit',
  'menu.ranges.list': 'Ranges List',
  'menu.ranges.edit': 'Ranges Edit',
  'menu.instances.list': 'Instances List',
  'menu.instances.edit': 'Instances Edit',
  'menu.excel.list': 'Excel List',
  'menu.excel.edit': 'Excel Edit',
  'menu.text.list': 'Text List',
  'menu.text.edit': 'Text Edit',
}

export default {
  ...components,
  ...locale
}
