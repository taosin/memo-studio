import React, { createContext, useState } from 'react';
import { IntlProvider } from 'react-intl';

import enLang from '../i18n/locales/en.json';
import zhLang from '../i18n/locales/zh.json';

const lang = {
	en: enLang,
	zh: zhLang,
}


// 创建上下文
export const LanguageContext = createContext({
	locale: 'en',
	setLocale: (locale: string) => {},
});

// 创建 Provider 组件
export const LanguageProvider: React.FC = ({ children }) => {
	const [locale, setLocale] = useState('zh'); // 默认语言

	return (
		<LanguageContext.Provider value={{ locale, setLocale }}>
			<IntlProvider locale={locale} messages={lang[locale]}>
				{children}
			</IntlProvider>
		</LanguageContext.Provider>
	);
};
