const sqlite3 = require('sqlite3').verbose();
const db = new sqlite3.Database('./db/notes.db');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');

class User {
  // 注册用户
  static register(username, password, callback) {
    const hashedPassword = bcrypt.hashSync(password, 10); // 加密密码
    db.run(
      'INSERT INTO users (username, password) VALUES (?, ?)',
      [username, hashedPassword],
      function (err) {
        if (err) return callback(err);
        callback(null, this.lastID);
      }
    );
  }

  // 登录用户
  static login(username, password, callback) {
    db.get('SELECT * FROM users WHERE username = ?', [username], (err, user) => {
      if (err) return callback(err);
      if (!user) return callback(new Error('用户不存在'));

      // 验证密码
      const isPasswordValid = bcrypt.compareSync(password, user.password);
      if (!isPasswordValid) return callback(new Error('密码错误'));

      // 生成 JWT
      const token = jwt.sign({
        id: user.id,
        username: user.username
      }, 'your-secret-key', {
        expiresIn: '1h', // Token 有效期
      });

      callback(null, token);
    });
  }
}

module.exports = User;
