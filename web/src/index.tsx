import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { ConfigProvider } from 'antd';
import { LanguageProvider } from './utils/LanguageContext';


// 获取根容器
const container = document.getElementById('root');
const root = createRoot(container!); // 创建根
root.render(
	<LanguageProvider>
		<ConfigProvider theme={{ token: { colorPrimary: '#00b96b' } }}>
			<App/>
		</ConfigProvider>
	</LanguageProvider>,
);
