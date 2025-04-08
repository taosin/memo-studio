// src/utils/storage.ts

interface StorageData<T> {
  value: T; // 存储的值
  expire: number; // 过期时间戳（毫秒）
}

class LocalStorage {
  /**
   * 设置缓存
   * @param key 键
   * @param value 值
   * @param expire 有效期（单位：秒），默认永不过期
   */
  static set<T>(key: string, value: T, expire?: number): void {
    const data: StorageData<T> = {
      value,
      expire: expire ? Date.now() + expire * 1000 : 0, // 0 表示永不过期
    };
    localStorage.setItem(key, JSON.stringify(data));
  }

  /**
   * 获取缓存
   * @param key 键
   * @returns 值或null（已过期/不存在）
   */
  static get<T>(key: string): T | null {
    const item = localStorage.getItem(key);
    if (!item) return null;

    try {
      const data = JSON.parse(item) as StorageData<T>;
      // 检查是否过期
      if (data.expire === 0 || data.expire >= Date.now()) {
        return data.value;
      }
      // 已过期则清除
      this.remove(key);
      return null;
    } catch (e) {
      // 解析失败则清除
      this.remove(key);
      return null;
    }
  }

  /**
   * 移除缓存
   * @param key 键
   */
  static remove(key: string): void {
    localStorage.removeItem(key);
  }

  /**
   * 清空所有缓存
   */
  static clear(): void {
    localStorage.clear();
  }

  /**
   * 检查键是否存在且未过期
   * @param key 键
   */
  static has(key: string): boolean {
    return this.get(key) !== null;
  }

  /**
   * 获取剩余有效期（秒）
   * @param key 键
   * @returns 剩余秒数（-1表示永不过期，-2表示不存在或已过期）
   */
  static getRemainTime(key: string): number {
    const item = localStorage.getItem(key);
    if (!item) return -2;

    try {
      const data = JSON.parse(item) as StorageData<unknown>;
      if (data.expire === 0) return -1;
      const remain = Math.ceil((data.expire - Date.now()) / 1000);
      return remain > 0 ? remain : -2;
    } catch {
      return -2;
    }
  }
}

export default LocalStorage;
