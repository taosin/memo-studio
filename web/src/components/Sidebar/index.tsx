import React from 'react';
import { FormattedMessage } from 'react-intl';
import { TagsOutlined } from '@ant-design/icons';
// @ts-ignore
import styles from './Sidebar.module.scss';
import TagTree from './TagTree';
import LanguageSwitcher from "../LanguageSwitcher";
import CalendarHeatmap from "../Heatmap";

const Sidebar: React.FC = () => {
	return (
		<div className={styles.sidebarContent}>
			{/* 日历区域 */}
			<div className={styles.calendarSection}>
				<h3><FormattedMessage id='calendar'/></h3>
				{/*<Calendar fullscreen={false}/>*/}
				<CalendarHeatmap/>
			</div>
			{/* 标签区域 */}
			<div className={styles.tagsSection}>
				<h3>
					<TagsOutlined/> <FormattedMessage id='tags'/>
				</h3>
				<div>
					<TagTree/>
				</div>
			</div>
			<LanguageSwitcher/>
		</div>
	);
};

export default Sidebar;
