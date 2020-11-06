import request from '../utils/request'

const api = {
  admin: '/admin',
  res: '/res',
  def: '/defs',
}

export default api

export function listDef () {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'listDef'}
  })
}

export function saveDef (data) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'createDef', 'data': data}
  })
}
