import {config} from '../utils/vari'

const PageSize = 15
const ResTypeDef = "def"
const ResTypeRanges = "ranges"
const ResTypeInstances = "instances"
const ResTypeConfig = "config"
export { PageSize, ResTypeDef, ResTypeRanges, ResTypeInstances, ResTypeConfig }

import {i18nRender} from '../locales'

export function checkLoop (rule, value, callback){
  console.log('checkLoop', value)

  const regx1 = /^[1-9][0-9]*$/;
  const regx2 = /^[1-9][0-9]*-?[1-9][0-9]*$/;
  if (!regx1.test(value) && !regx2.test(value)) {
    callback('需为整数或整数区间')
  }

  callback()
}

export function checkDirIsYaml (rule, value, callback){
  console.log('checkDirIsYaml', value)

  if (value.indexOf('yaml/') != 0 && value.indexOf('yaml\\') != 0) {
    callback(i18nRender('valid.folder.yaml'))
  }

  callback()
}
export function checkDirIsData (rule, value, callback){
  console.log('checkDirIsData', value)

  if (value.indexOf('data/') != 0 && value.indexOf('data\\') != 0) {
    callback(i18nRender('valid.folder.data'))
  }

  callback()
}
export function checkDirIsUsers (rule, value, callback){
  console.log('checkDirIsUsers', value)

  if (value.indexOf('users/') != 0 && value.indexOf('users\\') != 0) {
    callback(i18nRender('valid.folder.users'))
  }

  callback()
}

export function replacePathSep (path){
  const pth = path.replaceAll('\\', "/")
  return pth
}

export function pathToRelated (path){
  if (!config.workDir) return ''
  let name = path.substr(config.workDir.length)
  name = replacePathSep(name)
  return name
}

export function sectionStrToArr (str){
  str = str.substring(1, str.length - 1)
  let arr = str.split(',')
  str = arr.join('\n')
  return str
}

export function trimChar (str, ch){
  if (str.substr(0, 1) != ch || str.substr(str.length - 1, 1) != ch) {
    return str
  }

  if (str.indexOf(ch) == 0) {
    str = str.substring(1)
  }

  if (str.indexOf(ch) == str.length - 1) {
    str = str.substring(0, str.length - 1)
  }

  return str
}
