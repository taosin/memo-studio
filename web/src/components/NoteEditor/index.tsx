import React, { useState } from 'react';
// import ReactQuill from 'react-quill';
// import 'react-quill/dist/quill.snow.css';
// @ts-ignore
import styles from './NoteEditor.module.scss';

const NoteEditor: React.FC = () => {
	const [content, setContent] = useState<string>('');

	const handleChange = (value: string) => {
		setContent(value);
	};

	return (
		<div className={styles.editor}>
			{/*<ReactQuill*/}
			{/*	value={content}*/}
			{/*	onChange={handleChange}*/}
			{/*	placeholder="记录你的想法..."*/}
			{/*/>*/}
			dddd
		</div>
	);
};

export default NoteEditor;
