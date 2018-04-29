import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import ArrayInputs from '@/components/ArrayInputs'
import InputVersions from '@/components/InputVersions'
import ArrayList from '@/components/ArrayList'
import Vuetify from 'vuetify'

Vue.use(Vuetify)

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    },
    {
      path: '/accounts/:account_id/arrays/:array_id/inputs',
      name: 'ArrayInputs',
      component: ArrayInputs
    },
    {
      path: '/accounts/:account_id/arrays',
      name: 'ArrayList',
      component: ArrayList
    },
    {
      path: '/accounts/:account_id/arrays/:array_id/inputs/:input',
      name: 'InputVersions',
      component: InputVersions
    }
  ]
})
