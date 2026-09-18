import React from 'react';
import { useState } from 'react';
import { Editor } from '@monaco-editor/react';
import { GUInput } from '../ui/Input/GUInput';
import GUIButton from '../ui/Button/GUIButton';
import { useTranslation } from 'react-i18next';
import { useModalContext } from '@/context/useModalContext';

interface ModalAnalyticEditProps {
	title: string;
	query: string;
	onSave: (nextTitle: string, nextQuery: string) => void;
}

const normalizeQuery = (value: string) => {
	return (value ?? '')
		.replace(/\\r\\n/g, '\n')
		.replace(/\\n/g, '\n')
		.replace(/\\t/g, '\t')
		.trim();
};

const ModalAnalyticEdit: React.FC<ModalAnalyticEditProps> = ({ title, query, onSave }) => {
	const { t } = useTranslation();
	const { close } = useModalContext();

	const [draftTitle, setDraftTitle] = useState(title);
	const [draftQuery, setDraftQuery] = useState(normalizeQuery(query));

	return (
		<div className="space-y-3">
			<div>
				<label className="block text-sm mb-1">{t('label.title')}</label>
				<GUInput value={draftTitle} onChange={(e) => setDraftTitle(e.target.value)} className="w-full border border-[#d1d5db] rounded-md" />
			</div>

			<div>
				<label className="block text-sm mb-1">SQL</label>
				<Editor
					height={200}
					value={draftQuery}
					onChange={(val) => setDraftQuery(val ?? '')}
					defaultLanguage="sql"
					theme="vs-light"
					options={{
						minimap: { enabled: false },
						fontSize: 14,
						lineNumbers: 'on',
						lineNumbersMinChars: 3,
						scrollBeyondLastLine: false,
						wordWrap: 'off',

						scrollbar: {
							vertical: 'auto',
							horizontal: 'auto',
						},
					}}
				/>
			</div>

			<GUIButton
				type="submit"
				onClick={() => {
					onSave(draftTitle.trim(), normalizeQuery(draftQuery));
					close();
				}}
			>
				{t('label.save')}
			</GUIButton>
		</div>
	);
};

export default ModalAnalyticEdit;
