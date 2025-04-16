import React, { useState, useEffect, useRef } from 'react';
import { Input, message, Dropdown, Menu, List, Spin } from 'antd';
import { createNote } from '../../utils/api';
import SimpleMDE from 'simplemde';
import styles from './NoteEditor.module.scss';
import 'codemirror/addon/display/fullscreen';
import 'codemirror/lib/codemirror.css';
import 'codemirror/addon/display/fullscreen.css';

const { TextArea } = Input;

interface NoteEditorProps {
	onSave: () => void; // 保存后的回调函数，用于刷新笔记列表
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

	const editorRef = useRef<SimpleMDE | null>(null);
	const editorElementRef = useRef<HTMLTextAreaElement | null>(null);

	// 初始化 SimpleMDE 编辑器
	useEffect(() => {
		if (editorElementRef.current && !editorRef.current) {
			editorRef.current = new SimpleMDE({
				element: editorElementRef.current,
				spellChecker: false,
				status: false,
				toolbar: [
					'bold', 'italic', 'heading', '|',
					'quote', 'unordered-list', 'ordered-list', '|',
					'link', 'image', '|',
					'preview', 'side-by-side', 'fullscreen', '|',
					'guide'
				],
				initialValue: content,
				autofocus: true,
			});

			// 监听编辑器内容变化
			editorRef.current.codemirror.on('change', () => {
				const newContent = editorRef.current?.value() || '';
				setContent(newContent);

				// 获取光标位置
				const cursor = editorRef.current?.codemirror.getCursor();
				if (cursor) {
					const pos = editorRef.current?.codemirror.indexFromPos(cursor) || 0;
					setCursorPosition(pos);
				}

				// 检查是否输入了 # 或 @
				checkForSpecialCharacters(newContent, editorRef.current?.codemirror.indexFromPos(cursor) || 0);
			});

			// 监听光标位置变化
			editorRef.current.codemirror.on('cursorActivity', () => {
				const cursor = editorRef.current?.codemirror.getCursor();
				if (cursor) {
					const pos = editorRef.current?.codemirror.indexFromPos(cursor) || 0;
					setCursorPosition(pos);

					// 检查光标位置是否在 # 或 @ 后面
					const text = editorRef.current?.value() || '';
					checkForSpecialCharacters(text, pos);
				}
			});
		}

		return () => {
			if (editorRef.current) {
				editorRef.current.toTextArea();
				editorRef.current = null;
			}
		};
	}, []);

	// 当内容变化时更新编辑器
	useEffect(() => {
		if (editorRef.current && content !== editorRef.current.value()) {
			editorRef.current.value(content);
		}
	}, [content]);

	// 检查是否输入了特殊字符 (# 或 @)
	const checkForSpecialCharacters = (text: string, cursorPos: number) => {
		// 检查是否输入了 #
		const beforeCursor = text.substring(0, cursorPos);
		const lastHashIndex = beforeCursor.lastIndexOf('#');

		if (lastHashIndex !== -1 && lastHashIndex < cursorPos) {
			const searchText = beforeCursor.substring(lastHashIndex + 1, cursorPos);
			setTagSearchText(searchText);
			setShowTagDropdown(true);
			// 模拟获取标签列表
			fetchTags(searchText);
		} else {
			setShowTagDropdown(false);
		}

		// 检查是否输入了 @
		const lastAtIndex = beforeCursor.lastIndexOf('@');

		if (lastAtIndex !== -1 && lastAtIndex < cursorPos) {
			const searchText = beforeCursor.substring(lastAtIndex + 1, cursorPos);
			setMemoSearchText(searchText);
			setShowMemoDropdown(true);
			// 模拟获取备忘录列表
			fetchMemos(searchText);
		} else {
			setShowMemoDropdown(false);
		}
	};

	// 模拟获取标签列表
	const fetchTags = (searchText: string) => {
		setLoading(true);
		// 这里应该是实际的 API 调用
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
	};

	// 模拟获取备忘录列表
	const fetchMemos = (searchText: string) => {
		setLoading(true);
		// 这里应该是实际的 API 调用
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
	};

	// 插入标签
	const insertTag = (tag: Tag) => {
		if (editorRef.current) {
			const text = editorRef.current.value();
			const beforeCursor = text.substring(0, cursorPosition);
			const afterCursor = text.substring(cursorPosition);

			const lastHashIndex = beforeCursor.lastIndexOf('#');
			if (lastHashIndex !== -1) {
				const newText = beforeCursor.substring(0, lastHashIndex) + `#${tag.name} ` + afterCursor;
				editorRef.current.value(newText);
				setContent(newText);
				setShowTagDropdown(false);
			}
		}
	};

	// 插入备忘录引用
	const insertMemo = (memo: Memo) => {
		if (editorRef.current) {
			const text = editorRef.current.value();
			const beforeCursor = text.substring(0, cursorPosition);
			const afterCursor = text.substring(cursorPosition);

			const lastAtIndex = beforeCursor.lastIndexOf('@');
			if (lastAtIndex !== -1) {
				const newText = beforeCursor.substring(0, lastAtIndex) + `@${memo.title} ` + afterCursor;
				editorRef.current.value(newText);
				setContent(newText);
				setShowMemoDropdown(false);
			}
		}
	};

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
			if (editorRef.current) {
				editorRef.current.value('');
			}
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
				placeholder="标题"
				value={title}
				onChange={(e) => setTitle(e.target.value)}
				className={styles.titleInput}
			/>
			<div className={styles.editorContainer}>
				<textarea ref={editorElementRef} />

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
		</div>
	);
};

export default NoteEditor;
