import store from '../store'
import {
  APP_LANGUAGE,
} from '@/store/mutation-types'
import storage from 'store'

export default function Initializer () {
  store.dispatch('setLang', storage.get(APP_LANGUAGE, 'zh-CN'))
}
