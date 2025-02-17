import React from 'react';
import { Card, Button } from 'antd';
import styles from './NoteCard.module.scss';

interface NoteCardProps {
	title: string;
	content: string;
	date: string;
	onDelete: () => void;
}

const NoteCard: React.FC<NoteCardProps> = ({ title, content, date, onDelete }) => {
	return (
		<Card className={styles.noteCard} hoverable>
			<p>{content}</p>
			<small>{new Date(date).toLocaleString()}</small>
			<Button type="link" danger onClick={onDelete}>
				删除
			</Button>
		</Card>
	);
};

export default NoteCard;
