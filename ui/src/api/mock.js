import request from '../utils/request'

const mocksApi = '/mocks'

export function listMock (keywords, page) {
  return request({
    url: mocksApi,
    method: 'get',
    params: {keywords: keywords, page: page}
  })
}

export function getMock (id) {
  return request({
    url: `${mocksApi}/${id}`,
    method: 'get',
  })
}

export function saveMock (data) {
  return request({
    url: mocksApi,
    method: data.id ? 'put': 'post',
    data
  })
}

export function removeMock (id) {
  return request({
    url: `${mocksApi}/${id}`,
    method: 'delete'
  })
}

export function uploadMock (formData) {
  return request({
    url: `${mocksApi}/upload`,
    method: 'post',
    data: formData,
  })
}

export function getPreviewData (id) {
  return request({
    url: `${mocksApi}/getPreviewData`,
    method: 'get',
    params: {id},
  })
}

export function getPreviewResp (id, url, method, code, media) {
  return request({
    url: `${mocksApi}/getPreviewResp`,
    method: 'post',
    data: {id, url, method, code, media},
  })
}

export function startMockService (id, act) {
  return request({
    url: `${mocksApi}/startMockService`,
    method: 'post',
    params: {id, act}
  })
}
