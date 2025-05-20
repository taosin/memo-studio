import React, { useEffect, useState } from 'react';
import NoteCard from '../NoteCard';
import { fetchNotes, deleteNote, searchNotes } from '../../utils/api';
import styles from './NoteList.module.scss';
import SearchBar from "../SearchBar";
import NoteEditor from "../NoteEditor";

interface Note {
	id: number;
	title: string;
	content: string;
	created_at: string;
}

const NoteList: React.FC = () => {
	const [notes, setNotes] = useState<Note[]>([]);

	useEffect(() => {
		const loadNotes = async () => {
			const data = await fetchNotes();
			setNotes(data);
		};
		void loadNotes();
	}, []);

	const handleDelete = async (id: number) => {
		await deleteNote(id);
		setNotes(notes.filter((note) => note.id !== id));
	};

	const handleSearch = async (keyword: string) => {
		if (keyword.trim()) {
			const data = await searchNotes(keyword.trim());
			setNotes(data);
		} else {
			await loadNotes();
		}
	};

	// 加载笔记列表
	const loadNotes = async () => {
		const data = await fetchNotes();
		setNotes(data);
	};

	return (
		<div className={styles.noteList}>
			<SearchBar onSearch={handleSearch} />
			<NoteEditor onSave={loadNotes} />
			{notes.map((note) => (
				<NoteCard
					key={note.id}
					title={note.title}
					content={note.content}
					date={note.created_at}
					onDelete={() => handleDelete(note.id)}
				/>
			))}
		</div>
	);
};

export default NoteList;
