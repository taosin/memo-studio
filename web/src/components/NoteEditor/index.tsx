import React, { useEffect, useRef, useState } from 'react';
import SimpleMDE from 'simplemde';
import 'simplemde/dist/simplemde.min.css';
// @ts-ignore
import styles from './NoteEditor.module.scss';

const NoteEditor: React.FC = () => {
	const [value, setValue] = useState<string>('');
	const textareaRef = useRef<HTMLTextAreaElement>(null);
	const simpleMdeRef = useRef<SimpleMDE | null>(null);

	useEffect(() => {
		if (textareaRef.current) {
			simpleMdeRef.current = new SimpleMDE({
				element: textareaRef.current,
				initialValue: value,
				spellChecker: false, // 禁用拼写检查
				placeholder: '记录你的想法...',
				autoDownloadFontAwesome: false, // 禁用 FontAwesome 自动加载
			});

			simpleMdeRef.current.codemirror.on('change', () => {
				setValue(simpleMdeRef.current?.value() || '');
			});
		}

		return () => {
			simpleMdeRef.current?.toTextArea();
			simpleMdeRef.current = null;
		};
	}, []);

	return (
		<div className={styles.editor}>
			<textarea ref={textareaRef} />
		</div>
	);
};

export default NoteEditor;
