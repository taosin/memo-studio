import React from 'react';
import { createRoot } from 'react-dom/client';
import { IntlProvider } from 'react-intl';
import App from './App';
import { LanguageProvider } from './utils/LanguageContext';
import enLang from './i18n/locales/en.json';
import zhLang from './i18n/locales/zh.json';

const lang = {
	en: enLang,
	zh: zhLang,
}

const userLanguage = navigator.language.split('-')[0]; // 获取浏览器语言
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
