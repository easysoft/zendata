const PageSize = 15
const ResTypeDef = "def"
const ResTypeRanges = "ranges"
const ResTypeInstances = "instances"
export { PageSize, ResTypeDef, ResTypeRanges, ResTypeInstances }

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

  if (value.indexOf('yaml/') != 0) {
    callback('存放资源的目录必须以yaml/')
  }

  callback()
}
export function checkDirIsData (rule, value, callback){
  console.log('checkDirIsData', value)

  if (value.indexOf('data/') != 0) {
    callback('存放Excel的目录必须以data/')
  }

  callback()
}
export function checkDirIsUsers (rule, value, callback){
  console.log('checkDirIsUsers', value)

  if (value.indexOf('users/') != 0) {
    callback('存放数据的目录必须以users/')
  }

  callback()
}

