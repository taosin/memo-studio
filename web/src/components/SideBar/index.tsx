import React from 'react';
import { Calendar, Tag } from 'antd';
import { TagsOutlined } from '@ant-design/icons';
// @ts-ignore
import styles from './Sidebar.module.scss';

const Sidebar: React.FC = () => {
	return (
		<div className={styles.sidebarContent}>
			{/* 标签区域 */}
			<div className={styles.tagsSection}>
				<h3>
					<TagsOutlined /> 标签
				</h3>
				<div>
					<Tag color="magenta">工作</Tag>
					<Tag color="volcano">学习</Tag>
					<Tag color="orange">生活</Tag>
				</div>
			</div>

			{/* 日历区域 */}
			<div className={styles.calendarSection}>
				<h3>日历</h3>
				<Calendar fullscreen={false} />
			</div>
		</div>
	);
};

export default Sidebar;
