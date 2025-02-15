const express = require('express');
const Notes = require('../models/Note');
const router = express.Router();

// 获取所有笔记
router.get('/', (req, res) => {
  Notes.getAll((err, notes) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json(notes);
  });
});

// 创建笔记
router.post('/', (req, res) => {
  const { title, content } = req.body;
  if (!title || !content) {
    return res.status(400).json({ error: '标题和内容不能为空' });
  }

  Notes.create(title, content, (err, id) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json({ id, title, content });
  });
});

// 删除笔记
router.delete('/:id', (req, res) => {
  const { id } = req.params;
  Notes.delete(id, (err) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json({ success: true });
  });
});

// 搜索笔记
router.get('/search', (req, res) => {
  const { keyword } = req.query;
  if (!keyword) {
    return res.status(400).json({ error: '搜索关键词不能为空' });
  }

  Notes.search(keyword, (err, notes) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json(notes);
  });
});


module.exports = router;
