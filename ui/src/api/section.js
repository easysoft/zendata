import request from '../utils/request'
import api from './manage'

// section
export function listSection (ownerId, ownerType) {
  const data = {'action': 'listSection', id: ownerId, mode: ownerType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function createSection (ownerId, sectionId, ownerType) {
  const data = {'action': 'createSection',
      data: { ownerType: ownerType,  ownerId: ''+ownerId, sectionId: ''+sectionId}}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function updateSection (section, ownerType) {
  const data = {'action': 'updateSection', data: section, mode: ownerType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function removeSection (sectionId, ownerType) {
  const data = {'action': 'removeSection', id: sectionId, mode: ownerType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
