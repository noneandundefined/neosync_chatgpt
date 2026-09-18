import { useTranslation } from 'react-i18next';

interface EmptyStateTableProps {
	colSpan: number;
}

const EmptyStateTable: React.FC<EmptyStateTableProps> = ({ colSpan }) => {
	const { t } = useTranslation();

	return (
		<tr>
			<td colSpan={colSpan} id="td-data-available">
				<div className="flex justify-center my-[2rem] mt-[3rem]">
					<p className="text-[13px]">{t('message.data-available')}</p>
				</div>
			</td>
		</tr>
	);
};

export default EmptyStateTable;
