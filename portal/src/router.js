import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import ProblemSingle from './views/ProblemSingle.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/:domain',
      name: 'domain_home',
      component: Home,
      children: [
        {
          path: 'problem/:uid',
          name: 'problem_single',
          component: ProblemSingle,
        },]
    }
  ]
})
