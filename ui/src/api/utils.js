
export function checkLoop (rule, value, callback){
  console.log('checkLoop', value)

  const regx1 = /^[1-9][0-9]*$/;
  const regx2 = /^[1-9][0-9]*-?[1-9][0-9]*$/;
  if (!regx1.test(value) && !regx2.test(value)) {
    callback('需为整数或整数区间')
  }

  callback()
}
