import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import './plugins/element.js'
const moment = require('moment')
require('moment/locale/zh-cn')

import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.config.productionTip = false

Vue.use(VueAxios, axios)
Vue.use(require('vue-moment'), {
  moment
})
Vue.moment().locale('zh-cn')
Vue.use(require('tinymce'))

const EventBus = new Vue();

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
