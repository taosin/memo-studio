const express = require('express');
const Notes = require('../models/Note');
const router = express.Router();

// get all notes
router.get('/', (req, res) => {
  Notes.getAll((err, notes) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json(notes);
  });
});

// create note
router.post('/', (req, res) => {
  const { title, content } = req.body;
  if (!content) {
    return res.status(400).json({ error: 'title is null!' });
  }

  Notes.create(title, content, (err, id) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json({ id, title: title, content });
  });
});

// delete note by id
router.delete('/:id', (req, res) => {
  const { id } = req.params;
  Notes.delete(id, (err) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json({ success: true });
  });
});

// search notes by keyword
router.get('/search', (req, res) => {
  const { keyword } = req.query;
  if (!keyword) {
    return res.status(400).json({ error: 'key is null!' });
  }

  Notes.search(keyword, (err, notes) => {
    if (err) return res.status(500).json({ error: err.message });
    res.json(notes);
  });
});


module.exports = router;
