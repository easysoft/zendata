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

export function createMock (data) {
  return request({
    url: mocksApi,
    method: 'post',
    data
  })
}

export function updateMock (data) {
  return request({
    url: mocksApi,
    method: 'put',
    data
  })
}

export function removeMock (id) {
  return request({
    url: mocksApi,
    method: 'delete',
    params: {id: id}
  })
}
