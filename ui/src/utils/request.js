import notification from 'ant-design-vue/es/notification'
import axios from 'axios'
import { VueAxios } from './axios'

const request = axios.create({
  baseURL: getUrl(),
  timeout: 100000,
})

function getUrl() {
  let url = ''
  if (process.env.NODE_ENV === "development") {
    url = 'http://172.16.13.3:8848'
    console.log('dev env, url is ' + url)
  } else {
    const location = unescape(window.location.href);
    url = location.split('#')[0].split('index.html')[0];
    console.log('prod env, url is ' + url)
  }

  return url

}

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
  console.log('===Axios Request===', config.url, config.data);
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
