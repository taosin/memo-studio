import React, { useState } from 'react';
import { Input } from 'antd';
import { searchNotes } from '../../utils/api';
import styles from './index.module.scss';

const { Search } = Input;

interface SearchBarProps {
	onSearch: (notes: string) => void;
}

const SearchBar: React.FC<SearchBarProps> = ({ onSearch }) => {
	const [keyword, setKeyword] = useState('');

	const handleSearch = async () => {
		// const notes = await searchNotes(keyword);
		onSearch(keyword);
	};

	return (
		<div className={styles.searchBar}>
			<Search
				placeholder="输入关键词搜索笔记"
				value={keyword}
				onChange={(e) => setKeyword(e.target.value)}
				onSearch={handleSearch}
				enterButton
			/>
		</div>
	);
};

export default SearchBar;
