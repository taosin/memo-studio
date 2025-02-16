import React from "react";
import CalendarHeatmap from 'react-calendar-heatmap';
import 'react-calendar-heatmap/dist/styles.css';

const Heatmap: React.FC = () => {
	const getTooltipDataAttrs = (value) => {
		// Temporary hack around null value.date issue
		if (!value || !value.date) {
			return null;
		}
		// Configuration for react-tooltip
		return {
			'data-tip': `has count: ${value.count}`,
		};
	};
	const handleClick = (value) => {
		alert(`You clicked on ${value.date.toISOString().slice(0, 10)} with count: ${value.count}`);
	};

	return (
		<CalendarHeatmap
			startDate={new Date('2025-01-01')}
			endDate={new Date('2025-04-08')}
			values={[
				{ date: '2025-01-01', count: 12 },
				{ date: '2025-01-22', count: 122 },
				{ date: '2025-01-30', count: 38 },
				// ...and so on
			]}
			// showMonthLabels
			// showWeekdayLabels
			// classForValue={(value: string) => {
			// 	if (!value) {
			// 		return 'color-empty';
			// 	}
			// 	return `color-github-${value.count}`;
			// }}
			tooltipDataAttrs={getTooltipDataAttrs}
			onClick={handleClick}
		/>
	)
}
export default Heatmap;
