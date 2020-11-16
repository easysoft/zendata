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
export function getDef (id) {
  const data = {'action': 'getDef', id: id}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function saveDef (data) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'saveDef', 'data': data}
  })
}
export function saveDefDesign (data) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'saveDefDesign', 'data': data}
  })
}
export function removeDef (id) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'removeDef', id: id}
  })
}

export function getDefFieldTree (id) {
  const data = {'action': 'getDefFieldTree', id: id}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getDefField (id) {
  const data = {'action': 'getDefField', id: id}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function createDefField (targetId, mode) {
  const data = {'action': 'createDefField', id: targetId, mode: mode}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function saveDefField (data) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'saveDefField', 'data': data}
  })
}

export function removeDefField (id) {
  const data = {'action': 'removeDefField', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function moveDefField (src, dist, mode) {
  const data = {'action': 'moveDefField', src: src, dist: dist, mode: ''+mode}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listDefFieldSection (fieldId) {
  const data = {'action': 'listDefFieldSection', id: fieldId}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function createDefFieldSection (fieldId, sectionId) {
  const data = {'action': 'createDefFieldSection', data: { fieldId: ''+fieldId, sectionId: ''+sectionId}}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function updateDefFieldSection (section) {
  const data = {'action': 'updateDefFieldSection', data: section}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function removeDefFieldSection (sectionId) {
  const data = {'action': 'removeDefFieldSection', id: sectionId}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listDefFieldReferType (resType) {
  const data = {'action': 'listDefFieldReferType', mode: resType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getDefFieldRefer (fieldId, resType) {
  const data = {'action': 'getDefFieldRefer', id: fieldId, mode: resType}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function updateDefFieldRefer (refer) {
  const data = {'action': 'updateDefFieldRefer', data: refer}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
