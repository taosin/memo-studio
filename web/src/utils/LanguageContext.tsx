import React, { createContext, useState, useMemo, ReactNode, useEffect } from 'react';
import { IntlProvider } from 'react-intl';

import enLang from '../i18n/locales/en.json';
import zhLang from '../i18n/locales/zh.json';

const lang = {
	en: enLang,
	zh: zhLang,
};

// 默认语言常量
const DEFAULT_LOCALE = 'zh';

// 上下文类型定义
interface LanguageContextProps {
	locale: 'en' | 'zh';
	setLocale: (locale: 'en' | 'zh') => void;
}

// 创建上下文
export const LanguageContext = createContext<LanguageContextProps>({
	locale: DEFAULT_LOCALE,
	setLocale: () => { },
});

// 创建 Provider 组件
export const LanguageProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
	const [locale, setLocale] = useState<'en' | 'zh'>(() => {
		// 从本地缓存中获取语种
		return (localStorage.getItem('locale') as 'en' | 'zh') || DEFAULT_LOCALE;
	});

	// 使用 useMemo 优化上下文值
	const contextValue = useMemo(() => ({ locale, setLocale }), [locale]);

	// 当语种改变时，更新本地缓存
	useEffect(() => {
		localStorage.setItem('locale', locale);
	}, [locale]);

	return (
		<LanguageContext.Provider value={contextValue}>
			<IntlProvider locale={locale} messages={lang[locale]}>
				{children}
			</IntlProvider>
		</LanguageContext.Provider>
	);
};