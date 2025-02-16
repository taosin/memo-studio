import React from 'react';
import { ConfigProvider } from 'antd';
import AppLayout from './components/Layout';
import './styles/global.scss';

const App: React.FC = () => {
	return (
		<ConfigProvider theme={{ token: { colorPrimary: '#00b96b' } }}>
			<div className="app-container">
				<AppLayout/>
			</div>
		</ConfigProvider>
	)
};

export default App;
