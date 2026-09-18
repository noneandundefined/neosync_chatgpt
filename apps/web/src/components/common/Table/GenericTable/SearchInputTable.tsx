import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import Search from '@/components/@icons/search';
import { useTableContext } from './TableContext';

interface SearchInputTableProps {
	value: string;
	onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
	placeholder: string;
}

const SearchInputTable: React.FC<SearchInputTableProps> = ({ value, onChange, placeholder }) => {
	const { t } = useTranslation();

	const { setPage, setSearch } = useTableContext<any>();

	const handleSearch = () => {
		setSearch(value);
		setPage(1);
	};

	return (
		<form
			className="w-full flex justify-end gap-3"
			onSubmit={(e) => {
				e.preventDefault();
				handleSearch();
			}}
		>
			<input
				type="text"
				placeholder={placeholder}
				className="w-full md:min-w-[16rem] md:max-w-[35rem] bg-wgite border border-[#dedede] pl-3 pr-3 py-1.5 rounded text-sm"
				id="input__search__table"
				value={value}
				onChange={onChange}
			/>

			{value.length > 0 && (
				<Tooltip title={t('label.search')}>
					<button type="submit" id="buttonhlp" className="!p-[7px] md:!p-[5px] !h-auto">
						<Search fill="#49525f" size={20} />
					</button>
				</Tooltip>
			)}
		</form>
	);
};

export default SearchInputTable;
