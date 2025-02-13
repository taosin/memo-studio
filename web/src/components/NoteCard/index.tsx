import React from 'react';
import { Card } from 'antd';
import styles from './NoteCard.module.scss';

interface NoteCardProps {
	title: string;
	content: string;
	date: string;
}

const NoteCard: React.FC<NoteCardProps> = ({ title, content, date }) => {
	return (
		<Card className={styles.noteCard} hoverable>
			<h3>{title}</h3>
			<p>{content}</p>
			<small>{date}</small>
		</Card>
	);
};

export default NoteCard;
