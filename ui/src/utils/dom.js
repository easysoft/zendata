import config from '../config/config'
import { getElectron } from "./common";

export const setDocumentTitle = function (title) {
  document.title = title
  const ua = navigator.userAgent
  // eslint-disable-next-line
  const regex = /\bMicroMessenger\/([\d\.]+)/
  if (regex.test(ua) && /ip(hone|od|ad)/i.test(ua)) {
    const i = document.createElement('iframe')
    i.src = '/favicon.ico'
    i.style.display = 'none'
    i.onload = function () {
      setTimeout(function () {
        i.remove()
      }, 9)
    }
    document.body.appendChild(i)
  }
}

export const getPath = function () {
  return process.env.NODE_ENV === 'production' && !parseInt(process.env.UI_IN_CLIENT) && !getElectron() ? 'ui/': ''
}

export const domTitle = config.title
