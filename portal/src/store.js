import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    domain: {},
    user: {},
    guser: {},
    contest: {},
    contestTimer: {},
  },
  mutations: {
    setDomain(state, newDomain) {
      state.domain = JSON.parse(JSON.stringify(newDomain));
    },
    setUser(state, newUser) {
      state.user = JSON.parse(JSON.stringify(newUser));
    },
    setGUser(state, newUser) {
      state.guser = JSON.parse(JSON.stringify(newUser));
    },
    setContest(state, newContest) {
      state.contest = JSON.parse(JSON.stringify(newContest));
    },
    setContestTimer(state, newTimer) {
      state.contestTimer = JSON.parse(JSON.stringify(newTimer));
    }
  },
  actions: {

  }
})
