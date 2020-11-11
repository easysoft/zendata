import notification from 'ant-design-vue/es/notification'
import axios from 'axios'
import { VueAxios } from './axios'

const request = axios.create({
  baseURL: 'http://localhost:8848',
  timeout: 100000,
})

const errorHandler = error => {
  if (error.response) {
    const data = error.response.data

    if (error.response.status === 403) {
      notification.error({
        message: 'Forbidden',
        description: data.message
      })
    }
    if (error.response.status === 401) {
      notification.error({
        message: 'Unauthorized',
        description: 'Authorization verification failed'
      })
    }
  }
  return Promise.reject(error)
}

// request interceptor
request.interceptors.request.use(config => {
  return config
}, errorHandler)

// response interceptor
request.interceptors.response.use(response => {
  return response.data
}, errorHandler)

const installer = {
  vm: {},
  install (Vue) {
    Vue.use(VueAxios, request)
  }
}

export default request
export { installer as VueAxios, request as axios }
