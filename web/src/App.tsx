import React from 'react';
import AppLayout from './components/Layout';
import './styles/global.scss';

const App: React.FC = () => {
	return (
		<div className="app-container">
			<AppLayout/>
		</div>
	)
};

export default App;
