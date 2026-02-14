// ===== 加密工具 (Web Crypto API) =====

// 生成加密密钥
export async function generateEncryptionKey() {
  const key = await crypto.subtle.generateKey(
    { name: 'AES-GCM', length: 256 },
    true,
    ['encrypt', 'decrypt']
  );
  return key;
}

// 导出密钥为 base64
export async function exportKey(key) {
  const exported = await crypto.subtle.exportKey('raw', key);
  return btoa(String.fromCharCode(...new Uint8Array(exported));
}

// 从 base64 导入密钥
export async function importKey(base64Key) {
  const keyData = Uint8Array.from(atob(base64Key), c => c.charCodeAt(0));
  return await crypto.subtle.importKey(
    'raw',
    keyData,
    { name: 'AES-GCM', length: 256 },
    true,
    ['encrypt', 'decrypt']
  );
}

// 加密数据
export async function encryptData(plaintext, key) {
  const iv = crypto.getRandomValues(new Uint8Array(12));
  const encoder = new TextEncoder();
  const data = encoder.encode(plaintext);
  
  const encrypted = await crypto.subtle.encrypt(
    { name: 'AES-GCM', iv },
    key,
    data
  );
  
  const ivArray = Array.from(iv);
  const encryptedArray = Array.from(new Uint8Array(encrypted));
  return JSON.stringify({ iv: ivArray, data: encryptedArray });
}

// 解密数据
export async function decryptData(encrypted, key) {
  try {
    const { iv, data } = JSON.parse(encrypted);
    const ivArray = new Uint8Array(iv);
    const dataArray = new Uint8Array(data);
    
    const decrypted = await crypto.subtle.decrypt(
      { name: 'AES-GCM', iv: ivArray },
      key,
      dataArray
    );
    
    const decoder = new TextDecoder();
    return decoder.decode(decrypted);
  } catch (e) {
    console.error('解密失败:', e);
    return null;
  }
}

// 生成随机字符串
export function generateId(length = 32) {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  const randomValues = new Uint32Array(length);
  crypto.getRandomValues(randomValues);
  for (let i = 0; i < length; i++) {
    result += chars[randomValues[i] % chars.length];
  }
  return result;
}

// 本地存储加密封装
const ENCRYPTION_KEY_STORAGE_KEY = 'memo_encryption_key';

export async function getOrCreateEncryptionKey() {
  let key = localStorage.getItem(ENCRYPTION_KEY_STORAGE_KEY);
  if (!key) {
    const newKey = await generateEncryptionKey();
    key = await exportKey(newKey);
    localStorage.setItem(ENCRYPTION_KEY_STORAGE_KEY, key);
  }
  return importKey(key);
}

// 加密保存到 localStorage
export async function secureSave(key, data) {
  try {
    const encryptionKey = await getOrCreateEncryptionKey();
    const encrypted = await encryptData(JSON.stringify(data), encryptionKey);
    localStorage.setItem(key, encrypted);
    return true;
  } catch (e) {
    console.error('加密保存失败:', e);
    return false;
  }
}

// 从 localStorage 解密读取
export async function secureLoad(key) {
  try {
    const encrypted = localStorage.getItem(key);
    if (!encrypted) return null;
    
    const encryptionKey = await getOrCreateEncryptionKey();
    const decrypted = await decryptData(encrypted, encryptionKey);
    return decrypted ? JSON.parse(decrypted) : null;
  } catch (e) {
    console.error('解密读取失败:', e);
    return null;
  }
}

// 清除所有加密数据
export function clearSecureData() {
  localStorage.removeItem(ENCRYPTION_KEY_STORAGE_KEY);
}

// 隐私模式：检测是否在私密浏览模式下
export function isPrivateBrowsing() {
  // 尝试检测私密浏览
  try {
    const test = '__private_browser_test__';
    localStorage.setItem(test, test);
    localStorage.removeItem(test);
    return false;
  } catch (e) {
    return true;
  }
}

// 敏感数据脱敏
export function maskSensitiveData(data, fields = ['password', 'token', 'email', 'phone']) {
  if (typeof data !== 'object' || data === null) return data;
  
  const masked = { ...data };
  for (const field of fields) {
    if (masked[field]) {
      if (masked[field].length > 4) {
        masked[field] = masked[field].slice(0, 2) + '****' + masked[field].slice(-2);
      } else {
        masked[field] = '****';
      }
    }
  }
  return masked;
}
