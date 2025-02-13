import React from 'react';
import { Layout } from 'antd';
import Sidebar from '../Sidebar';
import NoteEditor from '../NoteEditor';
import NoteList from '../NoteList';
import styles from './Layout.module.scss';

const { Sider, Content } = Layout;

const AppLayout: React.FC = () => {
	return (
		<Layout className={styles.layout}>
			{/* 左侧边栏 */}
			<Sider width={320} theme="light" className={styles.sidebar}>
				<Sidebar />
			</Sider>

			{/* 右侧内容区域 */}
			<Layout>
				<Content className={styles.content}>
					<div className={styles.editorContainer}>
						{/*<NoteEditor />*/}
					</div>
					<div className={styles.notesContainer}>
						<NoteList />
					</div>
				</Content>
			</Layout>
		</Layout>
	);
};

export default AppLayout;
