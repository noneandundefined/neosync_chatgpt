import { useTranslation } from 'react-i18next';

interface CompactTableSearchProps {
	value: string;
	onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

const CompactTableSearch: React.FC<CompactTableSearchProps> = ({ value, onChange }) => {
	const { t } = useTranslation();

	return (
		<div className="w-full flex justify-end gap-3 mb-4">
			<input type="text" placeholder={t('label.search')} className="w-full md:min-w-[16rem] md:max-w-[35rem] bg-white border border-[#dedede] pl-3 pr-3 py-1.5 rounded text-sm" value={value} onChange={onChange} />
		</div>
	);
};

export default CompactTableSearch;
