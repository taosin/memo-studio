import React, { useState, useCallback } from 'react';
import { Input, message, List, Spin, Upload, Button } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
import { createNote } from '../../utils/api';
import styles from './NoteEditor.module.scss';

const { TextArea } = Input;

interface NoteEditorProps {
	onSave: () => void;
}

interface Tag {
	id: string;
	name: string;
}

interface Memo {
	id: string;
	title: string;
}

const NoteEditor: React.FC<NoteEditorProps> = ({ onSave }) => {
	const [title, setTitle] = useState<string>('');
	const [content, setContent] = useState<string>('');
	const [showTagDropdown, setShowTagDropdown] = useState<boolean>(false);
	const [showMemoDropdown, setShowMemoDropdown] = useState<boolean>(false);
	const [tagSearchText, setTagSearchText] = useState<string>('');
	const [memoSearchText, setMemoSearchText] = useState<string>('');
	const [tags, setTags] = useState<Tag[]>([]);
	const [memos, setMemos] = useState<Memo[]>([]);
	const [loading, setLoading] = useState<boolean>(false);
	const [cursorPosition, setCursorPosition] = useState<number>(0);

	// Quill 编辑器配置
	const modules = {
		toolbar: [
			[{ 'header': [1, 2, 3, 4, 5, 6, false] }],
			['bold', 'italic', 'underline', 'strike'],
			[{ 'list': 'ordered' }, { 'list': 'bullet' }],
			['link', 'image'],
			['clean']
		],
		clipboard: {
			matchVisual: false
		}
	};

	// 处理文件上传
	const handleFileUpload = useCallback(async (file: File) => {
		try {
			// 这里应该调用实际的文件上传 API
			const formData = new FormData();
			formData.append('file', file);

			// 模拟上传延迟
			await new Promise(resolve => setTimeout(resolve, 1000));

			// 返回文件 URL（这里使用模拟数据）
			const fileUrl = URL.createObjectURL(file);
			return fileUrl;
		} catch (error) {
			message.error('文件上传失败');
			return null;
		}
	}, []);

	// 检查特殊字符
	const checkForSpecialCharacters = useCallback((text: string, cursorPos: number) => {
		const beforeCursor = text.substring(0, cursorPos);

		// 检查标签
		const lastHashIndex = beforeCursor.lastIndexOf('#');
		if (lastHashIndex !== -1 && lastHashIndex < cursorPos) {
			const searchText = beforeCursor.substring(lastHashIndex + 1, cursorPos);
			setTagSearchText(searchText);
			setShowTagDropdown(true);
			fetchTags(searchText);
		} else {
			setShowTagDropdown(false);
		}

		// 检查备忘录引用
		const lastAtIndex = beforeCursor.lastIndexOf('@');
		if (lastAtIndex !== -1 && lastAtIndex < cursorPos) {
			const searchText = beforeCursor.substring(lastAtIndex + 1, cursorPos);
			setMemoSearchText(searchText);
			setShowMemoDropdown(true);
			fetchMemos(searchText);
		} else {
			setShowMemoDropdown(false);
		}
	}, []);

	// 获取标签列表
	const fetchTags = useCallback((searchText: string) => {
		setLoading(true);
		setTimeout(() => {
			const mockTags: Tag[] = [
				{ id: '1', name: '工作' },
				{ id: '2', name: '学习' },
				{ id: '3', name: '生活' },
				{ id: '4', name: '旅行' },
				{ id: '5', name: '健康' },
			].filter(tag => tag.name.includes(searchText));

			setTags(mockTags);
			setLoading(false);
		}, 300);
	}, []);

	// 获取备忘录列表
	const fetchMemos = useCallback((searchText: string) => {
		setLoading(true);
		setTimeout(() => {
			const mockMemos: Memo[] = [
				{ id: '1', title: '项目计划' },
				{ id: '2', title: '会议记录' },
				{ id: '3', title: '待办事项' },
				{ id: '4', title: '学习笔记' },
				{ id: '5', title: '旅行计划' },
			].filter(memo => memo.title.includes(searchText));

			setMemos(mockMemos);
			setLoading(false);
		}, 300);
	}, []);

	// 插入标签
	const insertTag = useCallback((tag: Tag) => {
		const quill = (document.querySelector('.ql-editor') as any)?.__quill;
		if (quill) {
			const range = quill.getSelection();
			if (range) {
				const text = quill.getText();
				const beforeCursor = text.substring(0, range.index);
				const lastHashIndex = beforeCursor.lastIndexOf('#');

				if (lastHashIndex !== -1) {
					quill.deleteText(lastHashIndex, range.index - lastHashIndex);
					quill.insertText(lastHashIndex, `#${tag.name} `);
				}
			}
		}
		setShowTagDropdown(false);
	}, []);

	// 插入备忘录引用
	const insertMemo = useCallback((memo: Memo) => {
		const quill = (document.querySelector('.ql-editor') as any)?.__quill;
		if (quill) {
			const range = quill.getSelection();
			if (range) {
				const text = quill.getText();
				const beforeCursor = text.substring(0, range.index);
				const lastAtIndex = beforeCursor.lastIndexOf('@');

				if (lastAtIndex !== -1) {
					quill.deleteText(lastAtIndex, range.index - lastAtIndex);
					quill.insertText(lastAtIndex, `@${memo.title} `);
				}
			}
		}
		setShowMemoDropdown(false);
	}, []);

	// 保存笔记
	const handleSave = async () => {
		if (!content.trim()) {
			message.warning('内容不能为空');
			return;
		}

		try {
			await createNote(title, content);
			message.success('笔记保存成功');
			setTitle('');
			setContent('');
			onSave();
		} catch (error) {
			message.error('保存失败，请重试');
		}
	};

	return (
		<div className={styles.editor}>
			<Input
				placeholder="标题"
				value={title}
				onChange={(e) => setTitle(e.target.value)}
				className={styles.titleInput}
			/>

			<div className={styles.editorContainer}>
				<ReactQuill
					value={content}
					onChange={(value) => {
						setContent(value);
						const quill = (document.querySelector('.ql-editor') as any)?.__quill;
						if (quill) {
							const range = quill.getSelection();
							if (range) {
								setCursorPosition(range.index);
								checkForSpecialCharacters(value, range.index);
							}
						}
					}}
					modules={modules}
					className={styles.quillEditor}
				/>

				{/* 标签下拉菜单 */}
				{showTagDropdown && (
					<div className={styles.dropdown}>
						<Spin spinning={loading}>
							<List
								dataSource={tags}
								renderItem={(tag) => (
									<List.Item
										className={styles.dropdownItem}
										onClick={() => insertTag(tag)}
									>
										#{tag.name}
									</List.Item>
								)}
							/>
						</Spin>
					</div>
				)}

				{/* 备忘录下拉菜单 */}
				{showMemoDropdown && (
					<div className={styles.dropdown}>
						<Spin spinning={loading}>
							<List
								dataSource={memos}
								renderItem={(memo) => (
									<List.Item
										className={styles.dropdownItem}
										onClick={() => insertMemo(memo)}
									>
										@{memo.title}
									</List.Item>
								)}
							/>
						</Spin>
					</div>
				)}
			</div>

			<div className={styles.toolbar}>
				<Upload
					customRequest={async ({ file }) => {
						const url = await handleFileUpload(file as File);
						if (url) {
							const quill = (document.querySelector('.ql-editor') as any)?.__quill;
							if (quill) {
								const range = quill.getSelection();
								if (range) {
									quill.insertEmbed(range.index, 'image', url);
								}
							}
						}
					}}
					showUploadList={false}
				>
					<Button icon={<UploadOutlined />}>上传图片</Button>
				</Upload>

				<Button type="primary" onClick={handleSave}>
					保存
				</Button>
			</div>
		</div>
	);
};

export default NoteEditor;
