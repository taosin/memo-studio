import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

const App = () => {
  const [notes, setNotes] = useState([]);
  const [newNote, setNewNote] = useState('');

  useEffect(() => {
    // fetchNotes();
  }, []);

  const fetchNotes = async () => {
    const response = await axios.get('/api/notes');
    setNotes(response.data);
  };

  const addNote = async () => {
    if (newNote.trim()) {
      const response = await axios.post('/api/notes', { content: newNote });
      setNotes([response.data, ...notes]);
      setNewNote('');
    }
  };

  return (
    <div className="app">
      <h1>Flomo Clone</h1>
      <div className="note-form">
        <textarea
          value={newNote}
          onChange={(e) => setNewNote(e.target.value)}
          placeholder="记录你的想法..."
        />
        <button onClick={addNote}>保存</button>
      </div>
      <div className="notes-list">
        {notes.map((note) => (
          <div key={note._id} className="note">
            <p>{note.content}</p>
            <small>{new Date(note.createdAt).toLocaleString()}</small>
          </div>
        ))}
      </div>
    </div>
  );
};

export default App;
