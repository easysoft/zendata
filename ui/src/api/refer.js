import request from '../utils/request'
import api from './manage'

export function getRefer (ownerId, ownerType) {
  const data = {'action': 'getRefer', id: ownerId, mode: ownerType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function updateRefer (refer, ownerType) {
  const data = {'action': 'updateRefer', data: refer, mode: ownerType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

// selection input on page
export function listReferType (resType) {
  const data = {'action': 'listReferType', mode: resType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function listReferField (refer) {
  const data = {'action': 'listReferField', data: refer}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
