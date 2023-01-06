
export const isWin = function () {
  const isWin = navigator.userAgentData.platform.toLowerCase().indexOf('win') > -1

  return isWin
}

export const getDir = function (name) {
  const ret = name + (isWin()? '\\' : '/')

  return ret
}
