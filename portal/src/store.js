import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    domain: {},
  },
  mutations: {
    setDomain(state, newDomain) {
      state.domain = JSON.parse(JSON.stringify(newDomain));
    }
  },
  actions: {

  }
})
