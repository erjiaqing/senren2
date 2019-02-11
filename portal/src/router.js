import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import ProblemSingle from './views/ProblemSingle.vue'
import ProblemList from './views/ProblemList.vue'
import ProblemEditor from './views/ProblemEditor.vue'
import ContestEditor from './views/ContestEditor.vue'
import ContestSingle from './views/ContestSingle.vue'
import ContestList from './views/ContestList.vue'
import SubmissionList from './views/SubmissionList.vue'
import SubmissionSingle from './views/SubmissionSingle.vue'
import DomainList from './views/DomainList.vue'
import DomainEdit from './views/DomainEditor.vue'
import DomainIndex from './views/DomainIndex.vue'
import Login from './views/Login.vue'

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
          path: "/",
          name: 'domain_index',
          component: DomainIndex,
        },
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
          path: 'contest/:uid',
          name: 'contest_single',
          component: ContestSingle,
        },
        {
          path: 'contest/:uid/edit',
          name: 'contest_editor',
          component: ContestEditor,
        },
        {
          path: 'contests',
          name: 'contest_list',
          component: ContestList,
        },
        {
          path: 'submissions/:filter',
          name: 'submission_list_filter',
          component: SubmissionList,
        },
        {
          path: 'submissions',
          name: 'submission_list',
          component: SubmissionList,
        },
        {
          path: 'submission/:uid',
          name: 'submission_single',
          component: SubmissionSingle,
        },
        {
          path: 'login',
          name: 'user_login',
          component: Login,
        },
        {
          path: 'domains',
          name: 'domain_list',
          component: DomainList,
        },
        {
          path: 'edit',
          name: 'domain_edit',
          component: DomainEdit,
        },
      ]
    }
  ]
})
