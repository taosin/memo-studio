import React from 'react';
import { createRoot } from 'react-dom/client';
import { IntlProvider } from 'react-intl';
import App from './App';

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
	<IntlProvider locale={userLanguage}
								messages={lang[userLanguage] || lang.en}>
		<App/>
	</IntlProvider>,
); // 使用新 API 渲染组件
// const root = ReactDOM.createRoot(document.getElementById('root'));
// root.render(
// 	<React.StrictMode>
// 		<App />
// 	</React.StrictMode>
// );
