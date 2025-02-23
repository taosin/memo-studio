import React from 'react';
import { Route, Redirect } from 'react-router-dom';

const PrivateRoute = ({ component: Component, ...rest }) => {
	const isAuthenticated = !!localStorage.getItem('token'); // 检查用户是否登录

	return (
		<Route
			{...rest}
			render={(props: React.JSX.IntrinsicAttributes) =>
				isAuthenticated ? (
					<Component {...props} />
				) : (
					<Redirect to="/login"/>
				)
			}
		/>
	);
};

export default PrivateRoute;
