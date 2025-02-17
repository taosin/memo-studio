import React, { useState } from "react";
import CalendarHeatmap from 'react-calendar-heatmap';
import dayjs from 'dayjs';
import 'react-calendar-heatmap/dist/styles.css';

const Heatmap: React.FC = () => {
	// 获取当前日期
	const currentDate = dayjs();

	// 计算本周的最后一天（周六）
	const endDate = currentDate.day(6); // 6 代表周六

	// 计算开始日期（结束日期减去72天）
	const startDate = endDate.subtract(84, 'day');

	const [hoveredDate, setHoveredDate] = useState(null);
	const [hoveredValue, setHoveredValue] = useState(null);

	const data = [
		{ date: '2025-01-01', count: 1 },
		{ date: '2025-01-02', count: 2 },
		{ date: '2025-01-23', count: 3 },
	];

	const [currentPosition, setCurrentPosition] = useState({
		top: 0,
		left: 0,
	});

	return (
		<>
			<CalendarHeatmap
				startDate={startDate}
				endDate={endDate}
				values={data}
				classForValue={(value: { count: any; }) => {
					if (!value) {
						return 'color-empty';
					}
					return `color-github-${value.count}`;
				}}
				onMouseOver={(event: any, value: {
					date: React.SetStateAction<null>;
					count: React.SetStateAction<null>;
				}) => {
					if (value) {
						console.error(event)
						setHoveredDate(value.date);
						setHoveredValue(value.count);
						setCurrentPosition({
							top: event.clientY + 10,
							left: event.clientX,
						})
					}
				}}
				onMouseLeave={() => {
					setHoveredDate(null);
					setHoveredValue(null);
				}}
				onClick={(value) => {
					console.log('Clicked date:', value?.date);
					console.log('Clicked count:', value?.count);
				}}
			/>
			{hoveredDate && (
				<div style={{
					position: 'fixed',
					background: '#fff',
					border: '1px solid #ccc',
					padding: '5px',
					width: '210px',
					top: currentPosition.top,
					left: currentPosition.left,
					zIndex: 99999
				}}>
					<span>{hoveredDate} 写了 {hoveredValue} 条笔记</span>
				</div>
			)}
		</>
	)
}
export default Heatmap;
