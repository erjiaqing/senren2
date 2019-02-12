import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    domain: {},
    user: {},
  },
  mutations: {
    setDomain(state, newDomain) {
      state.domain = JSON.parse(JSON.stringify(newDomain));
    },
    setUser(state, newUser) {
      state.user = JSON.parse(JSON.stringify(newUser));
    }
  },
  actions: {

  }
})
