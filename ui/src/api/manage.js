import request from '../utils/request'

const api = {
  admin: '',
  res: '/res',
  def: '/defs',
}

export default api

export function listDefs () {
  return request({
    url: api.admin,
    method: 'get',
    params: {}
  })
}

export function saveDef (data) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action' : 'createDef', 'data': data}
  })
}
