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
