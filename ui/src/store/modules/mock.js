import {CURR_MOCK_ITEM, CURR_DATA_SRC, CURR_MOCK_SRCS} from "@/store/mutation-types";
import {getMockDataSrc, getMockSataSrc, getPreviewData, listSampleSrc, saveMock} from "@/api/mock";

const mock = {
  state: {
    mockItem: {},
    mockSrcs: [],
    dataSrc: {}
  },
  mutations: {
    [CURR_MOCK_ITEM]: (state, item = {}) => {
      state.mockItem = item
    },
    [CURR_MOCK_SRCS]: (state, item = {}) => {
      state.mockSrcs = item
    },
    [CURR_DATA_SRC]: (state, data = {}) => {
      state.dataSrc = data
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

    previewMockItem ({ commit }, item) {
      return new Promise((resolve, reject) => {
        if (!item) {
          commit(CURR_MOCK_ITEM, null)
          resolve()
          return
        }

        listSampleSrc(item.id).then((json) => {
          commit(CURR_MOCK_SRCS, json.data)
        })

        getPreviewData(item.id).then((json) => {
          commit(CURR_MOCK_ITEM, json.data)

          const dataSrc = getMockDataSrc(json.data.item.paths)
          commit(CURR_DATA_SRC, dataSrc)

          resolve()
        }).catch(e => {
          reject(e)
        })

      })
    },
  }
}

export default mock
