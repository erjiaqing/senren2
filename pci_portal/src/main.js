import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import './plugins/element.js'
const moment = require('moment')

Vue.config.productionTip = false

const EventBus = new Vue();

Vue.use(require('vue-moment'), {
  moment
})
Vue.moment().locale('zh-cn')

Object.defineProperties(Vue.prototype, {
  $bus: {
    get: function () {
      return EventBus
    }
  }
});

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app')
