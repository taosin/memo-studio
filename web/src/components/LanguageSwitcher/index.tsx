import React, { useContext } from 'react';
import { Button } from 'antd';
import { LanguageContext } from '../../utils/LanguageContext';

const LanguageSwitcher: React.FC = () => {
	const { locale, setLocale } = useContext(LanguageContext);

	return (
		<div>
			<Button onClick={() => setLocale('en')} disabled={locale === 'en'}>
				English
			</Button>
			<Button onClick={() => setLocale('zh')} disabled={locale === 'zh'}>
				中文
			</Button>
		</div>
	);
};

export default LanguageSwitcher;
