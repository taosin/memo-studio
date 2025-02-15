const sqlite3 = require('sqlite3').verbose();
const db = new sqlite3.Database('./db/notes.db');

class Note {
  // 获取所有笔记
  static getAll(callback) {
    db.all('SELECT * FROM notes ORDER BY created_at DESC', callback);
  }

  // 创建笔记
  static create(title, content, callback) {
    db.run(
      'INSERT INTO notes (title, content) VALUES (?, ?)',
      [title, content],
      function (err) {
        if (err) return callback(err);
        callback(null, this.lastID);
      }
    );
  }

  // 删除笔记
  static delete(id, callback) {
    db.run('DELETE FROM notes WHERE id = ?', [id], callback);
  }

  // 搜索笔记
  static search(keyword, callback) {
    db.all(
      'SELECT * FROM notes WHERE title LIKE ? OR content LIKE ? ORDER BY created_at DESC',
      [`%${keyword}%`, `%${keyword}%`],
      callback
    );
  }
}

module.exports = Note;
