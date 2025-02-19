import React, { useState } from 'react';
import { Input, Button, message } from 'antd';
import { login } from '../../utils/api';
import { useHistory } from 'react-router-dom';
import styles from './index.module.scss';

const Index: React.FC = () => {
	const [username, setUsername] = useState('');
	const [password, setPassword] = useState('');
	const history = useHistory();

	const handleLogin = async () => {
		if (!username || !password) {
			message.warning('用户名和密码不能为空');
			return;
		}

		try {
			const { token } = await login(username, password);
			localStorage.setItem('token', token); // 保存 Token
			message.success('登录成功');
			history.push('/'); // 跳转到首页
		} catch (error) {
			message.error('登录失败，请检查用户名和密码');
		}
	};

	return (
		<div className={styles.login}>
			<h1>登录</h1>
			<Input
				placeholder="用户名"
				value={username}
				onChange={(e) => setUsername(e.target.value)}
				className={styles.input}
			/>
			<Input
				type="password"
				placeholder="密码"
				value={password}
				onChange={(e) => setPassword(e.target.value)}
				className={styles.input}
			/>
			<Button type="primary" onClick={handleLogin}>
				登录
			</Button>
		</div>
	);
};

export default Index;
