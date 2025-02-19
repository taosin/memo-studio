const express = require('express');
const User = require('../models/User');
const router = express.Router();

// 用户注册
router.post('/register', (req, res) => {
  const { username, password } = req.body;
  if (!username || !password) {
    return res.status(400).json({ error: '用户名和密码不能为空' });
  }

  User.register(username, password, (err, userId) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json({ id: userId, username });
  });
});

// 用户登录
router.post('/login', (req, res) => {
  const { username, password } = req.body;
  if (!username || !password) {
    return res.status(400).json({ error: '用户名和密码不能为空' });
  }

  User.login(username, password, (err, token) => {
    if (err) return res.status(401).json({ error: err.message });
    res.json({ token });
  });
});

module.exports = router;
