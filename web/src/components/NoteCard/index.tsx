import React from 'react';
import { Card, Button, Dropdown, MenuProps, Space } from 'antd';
import styles from './NoteCard.module.scss';
import { DownOutlined, SmileOutlined } from '@ant-design/icons';
import { FormattedMessage } from 'react-intl';

interface NoteCardProps {
	title: string;
	content: string;
	date: string;
	onDelete: () => void;
}

const NoteCard: React.FC<NoteCardProps> = ({ title, content, date, onDelete }) => {
	const items: MenuProps['items'] = [
		{

			key: '0',
			label: (
				< Button type="link" onClick={onDelete} >
					<FormattedMessage id='pin' />
					{/* 置顶 */}
				</Button >
			),
		},
		{
			key: '1',
			label: (
				< Button type="link" onClick={onDelete} >
					{/* 取消置顶 */}
					<FormattedMessage id='un-pin' />
				</Button >
			),
		},
		{
			key: '2',
			label: (
				< Button type="link" onClick={onDelete} >
					{/* 编辑 */}
					<FormattedMessage id='edit' />
				</Button >
			),
			icon: <SmileOutlined />,
			disabled: true,
		},
		{
			key: '3',
			label: (
				< Button type="link" danger onClick={onDelete} >
					{/* 删除 */}
					<FormattedMessage id='delete' />
				</Button >
			),
			disabled: true,
		},
	];
	return (
		<Card className={styles.noteCard}
			hoverable
			title={new Date(date).toLocaleString()}
			extra={<Dropdown menu={{ items }}>
				<a onClick={(e) => e.preventDefault()}>
					...
				</a>
			</Dropdown>}>
			{content}
		</Card>
	);
};

export default NoteCard;
