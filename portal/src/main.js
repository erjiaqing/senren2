import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import './plugins/element.js'

import axios from 'axios'
import VueAxios from 'vue-axios'

Vue.config.productionTip = false

Vue.use(VueAxios, axios)
Vue.use(require('vue-moment'))
Vue.use(require('tinymce'))

new Vue({
  router,
  store,
  render: function (h) { return h(App) }
}).$mount('#app')
