import request from '../utils/request'

const api = {
  admin: '/admin',
}

export default api

export function getWorkDir () {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'getWorkDir'}
  })
}

export function listDef (keywords, page) {
  return request({
    url: api.admin,
    method: 'post',
    data: {'action': 'listDef', keywords: keywords, page: page}
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

export function previewDefData (defId) {
  const data = {'action': 'previewDefData', id: defId}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function previewFieldData (fieldId, type) {
  const data = {'action': 'previewFieldData', id: fieldId, mode: type}
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listRanges (keywords, page) {
  const data = {'action': 'listRanges', keywords: keywords, page: page}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getRanges (id) {
  const data = {'action': 'getRanges', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveRanges (model) {
  const data = {'action': 'saveRanges', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeRanges (id) {
  const data = {'action': 'removeRanges', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function getResRangesItemTree (id) {
  const data = {'action': 'getResRangesItemTree', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getResRangesItem (id) {
  const data = {'action': 'getResRangesItem', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function createResRangesItem (rangesId, mode) {
  const data = {'action': 'createResRangesItem', domainId: rangesId, mode: mode}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveRangesItem (model) {
  const data = {'action': 'saveRangesItem', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}export function removeResRangesItem (itemId, rangesId) {
  const data = {'action': 'removeResRangesItem', id: itemId, domainId: rangesId}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listInstances (keywords, page) {
  const data = {'action': 'listInstances', keywords: keywords, page: page}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getInstances (id) {
  const data = {'action': 'getInstances', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveInstances (model) {
  const data = {'action': 'saveInstances', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeInstances (id) {
  const data = {'action': 'removeInstances', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function getResInstancesItemTree (id) {
  const data = {'action': 'getResInstancesItemTree', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getResInstancesItem (id) {
  const data = {'action': 'getResInstancesItem', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function createResInstancesItem (rangesId, mode) {
  const data = {'action': 'createResInstancesItem', domainId: rangesId, mode: mode}
  console.log(data)
  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveInstancesItem (model) {
  const data = {'action': 'saveInstancesItem', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeResInstancesItem (itemId, rangesId) {
  const data = {'action': 'removeResInstancesItem', id: itemId, domainId: rangesId}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function getResConfigItemTree (id) {
  const data = {'action': 'getResConfigItemTree', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getResConfigItem (id) {
  const data = {'action': 'getResConfigItem', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listText (keywords, page) {
  const data = {'action': 'listText', keywords: keywords, page: page}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getText (id) {
  const data = {'action': 'getText', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveText (model) {
  const data = {'action': 'saveText', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeText (id) {
  const data = {'action': 'removeText', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listConfig (keywords, page) {
  const data = {'action': 'listConfig', keywords: keywords, page: page}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getConfig (id) {
  const data = {'action': 'getConfig', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveConfig (model) {
  const data = {'action': 'saveConfig', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeConfig (id) {
  const data = {'action': 'removeConfig', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function listExcel (keywords, page) {
  const data = {'action': 'listExcel', keywords: keywords, page: page}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function getExcel (id) {
  const data = {'action': 'getExcel', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function saveExcel (model) {
  const data = {'action': 'saveExcel', data: model}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
export function removeExcel (id) {
  const data = {'action': 'removeExcel', id: id}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}

export function syncData () {
  const data = {'action': 'syncData', mode: ''}

  return request({
    url: api.admin,
    method: 'post',
    data: data
  })
}
