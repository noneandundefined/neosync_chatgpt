import React, { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

interface LabeledFieldProps {
	/* Ключ перевода для label (через i18n) */
	label: string;
	/* Дополнительный аргумент для форматирования текста label */
	argv?: string | number;
	/* Дополнительные CSS-классы для контейнера */
	className?: string;
	/* Дочерние элементы, которые будут отображаться внутри поля */
	children: ReactNode;
	/* Флаг блокировки поля */
	disabled?: boolean;
}

const LabeledField: React.FC<LabeledFieldProps> = ({ label, argv, className = 'flex flex-col gap-1', children, disabled = false }) => {
	const { t } = useTranslation();

	return (
		<div className={`${className} ${disabled ? 'opacity-50 pointer-events-none' : ''}`}>
			<label className="whitespace-nowrap text-sm">
				{t(label)} {argv}
			</label>
			{children}
		</div>
	);
};

export default LabeledField;
