import {APP_LANGUAGE, CURR_MOCK_ITEM} from "@/store/mutation-types";
import {loadLanguageAsync} from "@/locales";
import {previewMock, saveMock} from "@/api/mock";

const mock = {
  state: {
    mockItem: {}
  },
  mutations: {
    [CURR_MOCK_ITEM]: (state, item = {}) => {
      state.mockItem = item
    },
  },
  actions: {
    saveMockItem ({ commit }, item) {
      return new Promise((resolve, reject) => {
        saveMock(item).then((json) => {
            resolve(json)
          }).catch(e => {
            reject(e)
          })
      })
    },

    previewMockItem ({ commit }, id) {
      return new Promise((resolve, reject) => {
        previewMock(id).then((json) => {
          commit(CURR_MOCK_ITEM, json.data)
          resolve()
        }).catch(e => {
          reject(e)
        })
      })
    },
  }
}

export default mock
