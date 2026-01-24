// 认证状态管理
let token = null;
let user = null;

// 从 localStorage 读取
if (typeof window !== 'undefined') {
  token = localStorage.getItem('token');
  const savedUser = localStorage.getItem('user');
  if (savedUser) {
    try {
      user = JSON.parse(savedUser);
    } catch (e) {
      console.error('解析用户信息失败:', e);
    }
  }
}

const subscribers = new Set();

export const authStore = {
  subscribe(fn) {
    fn({ token, user, isAuthenticated: !!token });
    subscribers.add(fn);
    return () => subscribers.delete(fn);
  },
  setToken(newToken) {
    token = newToken;
    if (typeof window !== 'undefined') {
      if (newToken) {
        localStorage.setItem('token', newToken);
      } else {
        localStorage.removeItem('token');
      }
    }
    notify();
  },
  setUser(newUser) {
    user = newUser;
    if (typeof window !== 'undefined') {
      if (newUser) {
        localStorage.setItem('user', JSON.stringify(newUser));
      } else {
        localStorage.removeItem('user');
      }
    }
    notify();
  },
  login(newToken, newUser) {
    token = newToken;
    user = newUser;
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', newToken);
      localStorage.setItem('user', JSON.stringify(newUser));
    }
    notify();
  },
  logout() {
    token = null;
    user = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
    notify();
  },
  getToken() {
    return token;
  },
  getUser() {
    return user;
  },
  isAuthenticated() {
    return !!token;
  }
};

function notify() {
  subscribers.forEach(fn => fn({ token, user, isAuthenticated: !!token }));
}
