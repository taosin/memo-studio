import React, { useEffect, useState } from 'react';
import NoteCard from '../NoteCard';
import { fetchNotes, deleteNote } from '../../utils/api';
import styles from './NoteList.module.scss';

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
		loadNotes();
	}, []);

	const handleDelete = async (id: number) => {
		await deleteNote(id);
		setNotes(notes.filter((note) => note.id !== id));
	};

	return (
		<div className={styles.noteList}>
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
