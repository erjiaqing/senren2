import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import ProblemList from './views/ProblemList.vue'
import ProblemSingle from './views/ProblemSingle.vue'
import Login from './views/Login.vue'

Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/problems',
      name: 'problems',
      component: ProblemList
    },
    {
      path: '/problem/:uid',
      name: 'problem',
      component: ProblemSingle
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: function () { 
        return import(/* webpackChunkName: "about" */ './views/About.vue')
      }
    }
  ]
})
