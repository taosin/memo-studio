import React from 'react';
import NoteCard from '../NoteCard';
import styles from './NoteList.module.scss';

interface Note {
	id: number;
	title: string;
	content: string;
	date: string;
}

const notes: Note[] = [
	{
		id: 1,
		title: '工作笔记',
		content: '今天完成了项目需求分析...',
		date: '2023-10-01',
	},
	{
		id: 2,
		title: '学习笔记',
		content: '学习了 React Hooks 的使用...',
		date: '2023-10-02',
	},
	{
		id: 3,
		title: '生活笔记',
		content: '周末去公园散步，放松心情...',
		date: '2023-10-03',
	},
];

const NoteList: React.FC = () => {
	return (
		<div className={styles.noteList}>
			{notes.map((note) => (
				<NoteCard
					key={note.id}
					title={note.title}
					content={note.content}
					date={note.date}
				/>
			))}
		</div>
	);
};

export default NoteList;
