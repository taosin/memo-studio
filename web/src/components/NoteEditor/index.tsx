import React, { useState } from 'react';
import { Input, message } from 'antd';
import { createNote } from '../../utils/api';
import styles from './NoteEditor.module.scss';

const { TextArea } = Input;

interface NoteEditorProps {
	onSave: () => void; // 保存后的回调函数，用于刷新笔记列表
}

const NoteEditor: React.FC<NoteEditorProps> = ({ onSave }) => {
	const [title, setTitle] = useState<string>('');
	const [content, setContent] = useState<string>('');

	// 保存笔记
	const handleSave = async () => {
		if (!content.trim()) {
			message.warning('标题和内容不能为空');
			return;
		}

		try {
			await createNote(title, content);
			message.success('笔记保存成功');
			setTitle('');
			setContent('');
			onSave(); // 触发父组件刷新笔记列表
		} catch (error) {
			message.error('保存失败，请重试');
		}
	};

	// 监听键盘事件
	const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
		if (e.key === 'Enter' && e.shiftKey) {
			e.preventDefault(); // 阻止默认换行行为
			handleSave();
		}
	};

	return (
		<div className={styles.editor}>
			<Input
				placeholder="输入标题"
				value={title}
				onChange={(e) => setTitle(e.target.value)}
				className={styles.titleInput}
			/>
			<TextArea
				placeholder="记录你的想法...（Shift + Enter 保存）"
				value={content}
				onChange={(e) => setContent(e.target.value)}
				onKeyDown={handleKeyDown}
				autoSize={{ minRows: 4, maxRows: 8 }} // 自动调整高度
				className={styles.textarea}
			/>
		</div>
	);
};

export default NoteEditor;
