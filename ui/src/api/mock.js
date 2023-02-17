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

export function listSampleSrc (mockId) {
  return request({
    url: `${mocksApi}/${mockId}/listSampleSrc`,
    method: 'get'
  })
}

export function changeSampleSrc (mockId, key, value) {
  return request({
    url: `${mocksApi}/${mockId}/changeSampleSrc`,
    method: 'post',
    data: {key, value}
  })
}

export function getMockDataSrc (paths) {
  const dataSrc = {}

  Object.keys(paths).forEach((pathKey) => {
    const pathVal = paths[pathKey]

    Object.keys(pathVal).forEach((methodKey) => {
      const methodVal = pathVal[methodKey]

      Object.keys(methodVal).forEach((codeKey) => {
        const codeVal = methodVal[codeKey]

        Object.keys(codeVal).forEach((mediaKey) => {
          const samples = codeVal[mediaKey].samples

          const arr = ['schema']
          Object.keys(samples).forEach((sampleKey) => {
            // console.log(pathKey, methodKey, codeKey, mediaKey, sampleKey, samples[sampleKey])
            arr.push(sampleKey)
          })

          dataSrc[`${pathKey}-${methodKey}-${codeKey}-${mediaKey}`] = arr
        })

      })
    })
  })

  return dataSrc
}