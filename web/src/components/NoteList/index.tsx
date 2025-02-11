import React from 'react';
import { List, Avatar } from 'antd';
import { FileTextOutlined } from '@ant-design/icons';

interface Note {
	id: number;
	title: string;
	content: string;
}

const notes: Note[] = [
	{ id: 1, title: '工作笔记', content: '今天完成了项目需求分析...' },
	{ id: 2, title: '学习笔记', content: '学习了 React Hooks 的使用...' },
];

const NoteList: React.FC = () => {
	return (
		<List
			itemLayout="horizontal"
			dataSource={notes}
			renderItem={(note) => (
				<List.Item>
					<List.Item.Meta
						avatar={<Avatar icon={<FileTextOutlined/>}/>}
						title={note.title}
						description={note.content}
					/>
				</List.Item>
			)}
		/>
	);
};

export default NoteList;
