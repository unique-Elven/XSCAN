import { createRouter, createWebHistory } from 'vue-router'
import Layout from '../components/Layout.vue'
import Scanner from '../pages/Scanner.vue'
import Settings from '../pages/Settings.vue'
import History from '../pages/History.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: Layout,
      children: [
        { path: '', name: 'scanner', component: Scanner },
        { path: 'settings', name: 'settings', component: Settings },
        { path: 'history', name: 'history', component: History },
      ],
    },
  ],
})
