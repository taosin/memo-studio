import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import Home from '../pages/Home';
import Login from '../pages/Login';
import NotFound from '../pages/NotFound';

const Routers = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/home" element={<Home/>}/>
        <Route exact path="/" component={Home}/>

        {/* 登录页 */}
        <Route path="/login" component={Login}/>

        {/* 404 页面 */}
        <Route component={NotFound}/>
      </Routes>
    </BrowserRouter>
  );
};

export default Routers;
