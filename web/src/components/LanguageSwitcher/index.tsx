import React, { useContext } from 'react';
import { Button } from 'antd';
import { LanguageContext } from '../../utils/LanguageContext';
import styles from './LanguageSwitcher.module.scss';

const LanguageSwitcher: React.FC = () => {
	const { locale, setLocale } = useContext(LanguageContext);

	return (
		<div className={styles.switchLanguageContainer}>
			<div onClick={() => setLocale( locale === 'en' ? 'zh' : 'en')} className={styles.switchLanguageButton}>
				{
					locale === 'zh' ? '中文': 'English'
				}
			</div>
		</div>
	);
};

export default LanguageSwitcher;
