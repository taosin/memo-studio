const express = require('express');
const bodyParser = require('body-parser');
const cors = require('cors');
const authRouter = require('./routes/auth');
const notesRouter = require('./routes/notes');

const app = express();
const PORT = process.env.PORT || 9999;

// 中间件
app.use(cors());
app.use(bodyParser.json());

// API 路由
app.use('/api/auth', authRouter);
app.use('/api/notes', notesRouter);

// 启动服务器
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
