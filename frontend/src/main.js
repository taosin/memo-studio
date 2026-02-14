import App from './App.svelte';
import './styles/global.css';

const app = new App({
  target: document.getElementById('app'),
});

// 监听认证过期事件，提示用户重新登录
if (typeof window !== 'undefined') {
  window.addEventListener('auth-expired', () => {
    // 显示提示（如果有 toast 组件可以在这里调用）
    console.log('认证已过期，请重新登录');
    
    // 可选：自动跳转到登录页
    // window.location.href = '/login';
  });
}

export default app;
