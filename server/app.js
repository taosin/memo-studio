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

// Start the server, but not when required as a module (for testing).
if (require.main === module) {
  app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
  });
}

module.exports = app;
