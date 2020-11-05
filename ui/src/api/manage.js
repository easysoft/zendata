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

export function saveDef () {
  return request({
    url: api.admin,
    method: 'post',
    data: {'name' : 123}
  })
}
