import notification from 'ant-design-vue/es/notification'
import axios, {AxiosInstance} from 'axios'

let request = null
initRequest()

// used to switch to another remote service
function initRequest(remoteUrl) {
  const url = remoteUrl ? remoteUrl : getUrl()

  request = axios.create({
    baseURL: url,
    timeout: 100000,
  })
}

function getUrl() {
  let url = ''
  if (process.env.NODE_ENV === "development") {
    url = 'http://127.0.0.1:8848/api/v1'
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
  console.log('---Request---', config.url);
  return config
}, errorHandler)

// response interceptor
request.interceptors.response.use(resp => {
  console.log('---Response---', resp.config.url, resp.data);
  return resp.data
}, errorHandler)

export default request
