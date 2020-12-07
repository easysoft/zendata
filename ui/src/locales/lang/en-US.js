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

  'title.data.create': 'Data Creation',
  'title.config.create': 'Config Creation',
  'title.ranges.create': 'Ranges Creation',
  'title.instances.create': 'Instances Creation',
  'title.excel.create': 'Excel Creation',
  'title.text.create': 'Text Creation',

  'msg.mine': 'My Data',
  'msg.buildin': 'Build In',

  'msg.workdir': 'WorkDir',
  'msg.help': 'Help',
  'msg.yes': 'Yes',
  'msg.no': 'No',

  'msg.data': 'Data',
  'msg.config': 'Config',
  'msg.ranges': 'Ranges',
  'msg.instances': 'Instances',
  'msg.excel': 'Excel',
  'msg.text': 'Text',

  'action.list': 'List',
  'action.create': 'Create',
  'action.edit': 'Edit',
  'action.delete': 'Delete',
  'action.preview': 'Preview',
  'action.design': 'Design',
  'action.back': 'Back',
  'action.save': 'Save',
  'action.reset': 'Reset',
  'action.import.from.file': 'Import From Files',

  'tips.refer': 'Reference',
  'tips.pls.select': 'Please select',
  'tips.range.int': 'Integer or a range of integers',

  'tips.delete': 'Are you sure to delete?',
  'tips.search': 'Input keywords to search',
  'tips.success.to.import': 'Success to import.',

  'form.name': 'Name',
  'form.file': 'File',
  'form.dir': 'Directory',
  'form.sub.dir': 'Sub Directory',
  'form.folder': 'Folder',
  'form.path': 'Path',
  'form.file.name': 'File Name',
  'form.desc': 'Description',
  'form.content': 'Content',
  'form.prefix': 'Prefix',
  'form.postfix': 'Postfix',
  'form.loop': 'Loop',
  'form.loopfix': 'LoopFix',
  'form.format': 'Format',
  'form.function': 'Function',
  'form.opt': 'Operation',

  'valid.required': 'Can not be empty.',
  'valid.loop.format': 'Should be integer or a range of integers',
  'valid.folder.users': 'Data must be saved in users/ dir',
  'valid.folder.yaml': 'YAML must be saved in /yaml',
  'valid.folder.data': 'Excel must be saved in data/ dir',

}

export default {
  ...components,
  ...locale
}
