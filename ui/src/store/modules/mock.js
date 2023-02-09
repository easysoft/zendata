import {APP_LANGUAGE, CURR_MOCK_ITEM} from "@/store/mutation-types";
import {loadLanguageAsync} from "@/locales";
import {getPreviewData, previewMock, saveMock} from "@/api/mock";

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
    setMockItem ({ commit }, item) {
      return new Promise((resolve, reject) => {
        commit(CURR_MOCK_ITEM, item)
        resolve()
      })
    },

    saveMockItem ({ commit }, item) {
      return new Promise((resolve, reject) => {
        saveMock(item).then((json) => {
            resolve(json)
          }).catch(e => {
            reject(e)
          })
      })
    },

    previewMockItem ({ commit }, item) {
      return new Promise((resolve, reject) => {
        if (!item) {
          commit(CURR_MOCK_ITEM, null)
          resolve()
          return
        }

        getPreviewData(item.id).then((json) => {
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
