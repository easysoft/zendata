import Vue from 'vue'
import Vuex from 'vuex'

import app from './modules/app'
import mock from './modules/mock'
import getters from './getters'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    app,mock
  },
  state: {

  },
  mutations: {

  },
  actions: {

  },
  getters
})
