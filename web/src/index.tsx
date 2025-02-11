import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';

// 获取根容器
const container = document.getElementById('root');
const root = createRoot(container!); // 创建根
root.render(<App />); // 使用新 API 渲染组件
// const root = ReactDOM.createRoot(document.getElementById('root'));
// root.render(
// 	<React.StrictMode>
// 		<App />
// 	</React.StrictMode>
// );
