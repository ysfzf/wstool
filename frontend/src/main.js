import {createApp} from 'vue'
import App from './App.vue'
import './style.css';
import '@layui/layui-vue/lib/index.css'
import Layui from '@layui/layui-vue'
createApp(App).use(Layui).mount('#app')
