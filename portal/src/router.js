import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import ProblemSingle from './views/ProblemSingle.vue'
import ProblemList from './views/ProblemList.vue'
import ProblemEditor from './views/ProblemEditor.vue'
import SubmissionSingle from './views/SubmissionSingle.vue'

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
          path: 'problems',
          name: 'problem_list',
          component: ProblemList,
        },
        {
          path: 'problem/:uid',
          name: 'problem_single',
          component: ProblemSingle,
        },
        {
          path: 'problem/:uid/edit',
          name: 'problem_editor',
          component: ProblemEditor,
        },
        {
          path: 'submission/:uid',
          name: 'submission_single',
          component: SubmissionSingle,
        },
      ]
    }
  ]
})
