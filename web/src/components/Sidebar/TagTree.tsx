import React from 'react';
import { Tree } from 'antd';
import { FolderOutlined, FileOutlined } from '@ant-design/icons';
import type { DataNode } from 'antd/es/tree';
import styles from './TagTree.module.scss';

const { DirectoryTree } = Tree;

// 树形标签数据
const treeData: DataNode[] = [
	{
		title: '工作',
		key: 'work',
		icon: <FolderOutlined />,
		children: [
			{
				title: '项目 A',
				key: 'work-project-a',
				icon: <FileOutlined />,
			},
			{
				title: '项目 B',
				key: 'work-project-b',
				icon: <FileOutlined />,
			},
		],
	},
	{
		title: '学习',
		key: 'study',
		icon: <FolderOutlined />,
		children: [
			{
				title: 'React',
				key: 'study-react',
				icon: <FileOutlined />,
			},
			{
				title: 'TypeScript',
				key: 'study-typescript',
				icon: <FileOutlined />,
			},
		],
	},
	{
		title: '生活',
		key: 'life',
		icon: <FolderOutlined />,
		children: [
			{
				title: '旅行计划',
				key: 'life-travel',
				icon: <FileOutlined />,
			},
			{
				title: '购物清单',
				key: 'life-shopping',
				icon: <FileOutlined />,
			},
		],
	},
];

const TagTree: React.FC = () => {
	// 处理节点点击事件
	const handleSelect = (selectedKeys: React.Key[], info: any) => {
		console.log('Selected:', selectedKeys, info);
	};

	return (
		<div className={styles.tagTree}>
			<DirectoryTree
				defaultExpandAll
				onSelect={handleSelect}
				treeData={treeData}
			/>
		</div>
	);
};

export default TagTree;
