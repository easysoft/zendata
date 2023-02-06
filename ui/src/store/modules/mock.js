import {CURR_MOCK_ITEM} from "@/store/mutation-types";

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
      return new Promise((resolve) => {
        commit(CURR_MOCK_ITEM, item)
        resolve()
      })
    }
  }
}

export default mock
