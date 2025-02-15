import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { LanguageProvider } from './utils/LanguageContext';


// 获取根容器
const container = document.getElementById('root');
const root = createRoot(container!); // 创建根
root.render(
	<LanguageProvider>
		<App/>
	</LanguageProvider>,
); // 使用新 API 渲染组件
// const root = ReactDOM.createRoot(document.getElementById('root'));
// root.render(
// 	<React.StrictMode>
// 		<App />
// 	</React.StrictMode>
// );
