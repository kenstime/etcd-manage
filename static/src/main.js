import Vue from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import './plugins/iview.js'

import VueCodeMirror from 'vue-codemirror-lite'

Vue.use(VueCodeMirror)


Vue.config.productionTip = false;

axios.defaults.baseURL = 'http://127.0.0.1:9090';

// 请求拦截器

axios.interceptors.request.use(function (config) {
  // 读取本地是否存在用户信息
  let etcdUserStr = localStorage.getItem('etcdUser');
  console.log(etcdUserStr)
  if(etcdUserStr != null && typeof etcdUserStr != 'undefined'){
    let etcdUser = JSON.parse(etcdUserStr);
    if (typeof etcdUser.username != 'undefined' && etcdUser.username != ''){
      config.headers['X-Etcd-Username'] = etcdUser.username;
    }
    if (typeof etcdUser.password != 'undefined' && etcdUser.password != ''){
      config.headers['X-Etcd-Password'] = etcdUser.password;
    }
  }
  return config;
}, function (error) {
  return Promise.reject(error);
})


Vue.prototype.$http = axios;

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
