const express = require('express');
const mongoose = require('mongoose');
const bodyParser = require('body-parser');
const cors = require('cors');

const app = express();
const PORT = process.env.PORT || 9999;

// 中间件
app.use(cors());
app.use(bodyParser.json());

// 连接 MongoDB
mongoose.connect('mongodb://localhost:27017/flomo', {
  useNewUrlParser: true,
  useUnifiedTopology: true,
});

// 定义笔记模型
const NoteSchema = new mongoose.Schema({
  content: String,
  createdAt: { type: Date, default: Date.now },
});

const Note = mongoose.model('Note', NoteSchema);

// API 路由
app.get('/api/notes', async (req, res) => {
  const notes = await Note.find().sort({ createdAt: -1 });
  res.json(notes);
});

app.post('/api/notes', async (req, res) => {
  const newNote = new Note({ content: req.body.content });
  await newNote.save();
  res.json(newNote);
});

app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
