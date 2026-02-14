import { secureSave, secureLoad, getOrCreateEncryptionKey, generateId } from './encryption.js';

// 自动保存配置
const AUTO_SAVE_KEY = 'memo_autosave_draft';
const MAX_DRAFTS = 10;
const AUTO_SAVE_INTERVAL = 30000; // 30秒

let autoSaveTimer = null;
let currentDraft = null;

// 保存草稿
export async function saveDraft(note) {
  try {
    const drafts = await getDrafts();
    
    // 更新或添加
    const existingIndex = drafts.findIndex(d => d.id === note.id);
    if (existingIndex >= 0) {
      drafts[existingIndex] = {
        ...note,
        updatedAt: new Date().toISOString(),
        draftId: drafts[existingIndex].draftId
      };
    } else {
      drafts.unshift({
        ...note,
        draftId: generateId(16),
        updatedAt: new Date().toISOString()
      });
    }
    
    // 限制数量
    if (drafts.length > MAX_DRAFTS) {
      drafts.splice(MAX_DRAFTS);
    }
    
    await secureSave(AUTO_SAVE_KEY, drafts);
    return true;
  } catch (e) {
    console.error('保存草稿失败:', e);
    return false;
  }
}

// 获取所有草稿
export async function getDrafts() {
  const drafts = await secureLoad(AUTO_SAVE_KEY);
  return Array.isArray(drafts) ? drafts : [];
}

// 获取单个草稿
export async function getDraft(noteId) {
  const drafts = await getDrafts();
  return drafts.find(d => d.id === noteId || d.draftId === noteId);
}

// 删除草稿
export async function deleteDraft(draftId) {
  const drafts = await getDrafts();
  const filtered = drafts.filter(d => d.draftId !== draftId);
  await secureSave(AUTO_SAVE_KEY, filtered);
}

// 清除所有草稿
export async function clearDrafts() {
  localStorage.removeItem(AUTO_SAVE_KEY);
}

// 开始自动保存
export function startAutoSave(note, onSave) {
  if (autoSaveTimer) {
    clearInterval(autoSaveTimer);
  }
  
  currentDraft = note;
  autoSaveTimer = setInterval(async () => {
    if (currentDraft) {
      await saveDraft(currentDraft);
      if (onSave) onSave();
    }
  }, AUTO_SAVE_INTERVAL);
}

// 停止自动保存
export function stopAutoSave() {
  if (autoSaveTimer) {
    clearInterval(autoSaveTimer);
    autoSaveTimer = null;
  }
  currentDraft = null;
}

// ===== 数据备份 =====

const BACKUP_KEY = 'memo_backup_list';

// 创建备份
export async function createBackup(notes, tags) {
  const backup = {
    id: generateId(32),
    createdAt: new Date().toISOString(),
    version: '1.0',
    data: {
      notes: notes.map(n => ({
        ...n,
        // 不保存标签详情，只保存ID
        tags: n.tags?.map(t => ({ id: t.id, name: t.name, color: t.color })) || []
      })),
      tags: tags
    }
  };
  
  // 保存备份索引
  const backups = await getBackupList();
  backups.unshift({
    id: backup.id,
    createdAt: backup.createdAt,
    noteCount: notes.length,
    tagCount: tags.length
  });
  
  // 只保留最近5个备份
  if (backups.length > 5) {
    backups.length = 5;
  }
  
  localStorage.setItem(BACKUP_KEY, JSON.stringify(backups));
  
  // 保存完整备份数据
  const backupDataKey = `memo_backup_${backup.id}`;
  localStorage.setItem(backupDataKey, JSON.stringify(backup));
  
  return backup;
}

// 获取备份列表
export function getBackupList() {
  try {
    return JSON.parse(localStorage.getItem(BACKUP_KEY)) || [];
  } catch {
    return [];
  }
}

// 获取完整备份
export function getBackup(backupId) {
  try {
    return JSON.parse(localStorage.getItem(`memo_backup_${backupId}`));
  } catch {
    return null;
  }
}

// 删除备份
export function deleteBackup(backupId) {
  localStorage.removeItem(`memo_backup_${backupId}`);
  const backups = getBackupList().filter(b => b.id !== backupId);
  localStorage.setItem(BACKUP_KEY, JSON.stringify(backups));
}

// 导出数据为 JSON 文件
export function exportAsJSON(notes, tags) {
  const data = {
    exportDate: new Date().toISOString(),
    version: '1.0',
    notes,
    tags
  };
  
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `memo-studio-backup-${new Date().toISOString().split('T')[0]}.json`;
  a.click();
  URL.revokeObjectURL(url);
}

// 导入 JSON 文件
export function importFromJSON(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = async (e) => {
      try {
        const data = JSON.parse(e.target.result);
        resolve(data);
      } catch (err) {
        reject(new Error('无效的备份文件'));
      }
    };
    reader.onerror = () => reject(new Error('读取文件失败'));
    reader.readAsText(file);
  });
}

// ===== 离线支持 =====

// 检测离线状态
export function isOnline() {
  return navigator.onLine;
}

// 注册离线事件监听
export function registerOfflineListener(callback) {
  window.addEventListener('offline', () => callback(false));
  window.addEventListener('online', () => callback(true));
}

// 服务工作者注册
export async function registerServiceWorker() {
  if ('serviceWorker' in navigator) {
    try {
      const registration = await navigator.serviceWorker.register('/sw.js');
      console.log('Service Worker 注册成功:', registration.scope);
      return registration;
    } catch (e) {
      console.error('Service Worker 注册失败:', e);
      return null;
    }
  }
  return null;
}
