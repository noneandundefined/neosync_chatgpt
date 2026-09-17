import i18next from 'i18next';
import { toast } from 'react-toastify';
import FileHelper from '@/utils/FilesUtils';
import Tooltip, { PositionType } from '@/components/ui/Tooltip';

interface ActionButtonProps {
	/**
	 * Иконка, которая будет отображаться в кнопке.
	 * Обычно это React-компонент из библиотеки (например lucide-react).
	 */
	icon: any;
	/** Подпись под иконкой. */
	label: string;
	/**
	 * Флаг, блокирующий кнопку.
	 * Если true – кнопка недоступна для клика.
	 */
	disabled?: boolean;
	/**
	 * Флаг изменения состояния.
	 * Если true – кнопка подсвечивается зелёным цветом.
	 */
	isChange?: boolean | null;
	/** Колбэк, вызываемый при клике на кнопку. */
	onClick?: () => void;
	/** Если это кнопка выбора файла */
	isFileInput?: boolean;
	onFileSelect?: (file: File) => void;
	tooltip?: string;
	tooltipPosition?: PositionType;
}

const ActionButton: React.FC<ActionButtonProps> = ({ icon: Icon, label, disabled, isChange = false, onClick, isFileInput, onFileSelect, tooltip, tooltipPosition = 'top' }) => {
	const handleClick = (e: React.MouseEvent) => {
		if (disabled) {
			e.stopPropagation();
			return;
		}

		onClick?.();
	};

	const button = isFileInput ? (
		<label className={`max-w-[5rem] lg:max-w-[4rem] flex flex-col items-center justify-start text-black ${disabled ? 'cursor-not-allowed' : 'cursor-pointer'}`}>
			<input
				type="file"
				className="hidden"
				disabled={disabled}
				onChange={async (e) => {
					if (!e.target.files || !e.target.files[0] || !onFileSelect) return;

					const file = e.target.files[0];

					const isBinary = await FileHelper.isBin(file);
					if (!isBinary) {
						toast.error(i18next.t('message.file-not-binary'));
						return;
					}

					if (!FileHelper.fileMax5Kb(file)) {
						toast.error(i18next.t('message.file-too-large-5kb'));
						return;
					}

					await onFileSelect(file);

					e.target.value = '';
				}}
			/>
			<Icon fill="#000" size={41} className={`border-2 border-[#d1d5db] rounded-[8px] active:scale-90 bg-[#f3f3f3] ${disabled ? 'opacity-50' : 'hover:bg-[#f9f9f9]'} p-2`} />
			<p className="text-[12px] text-center mt-2">{label}</p>
		</label>
	) : (
		<div className={`max-w-[4rem] flex flex-col items-center justify-start ${disabled ? 'cursor-not-allowed' : 'cursor-pointer'}`} onClick={handleClick}>
			<Icon
				fill={disabled ? '#9ca3af' : '#000'}
				size={41}
				className={`border-2 rounded-[8px] active:scale-90 ${
					disabled ? 'border-[#e5e7eb] bg-[#f3f3f3]' : isChange ? 'border-[#7db500] bg-[#7db500] hover:bg-[#94d505]' : 'border-[#d1d5db] bg-[#f3f3f3] hover:bg-[#f9f9f9]'
				} p-2`}
			/>
			<p className={`text-[12px] text-center mt-2 ${disabled ? 'text-[#9ca3af]' : ''}`}>{label}</p>
		</div>
	);

	if (!tooltip) {
		return button;
	}

	return (
		<Tooltip title={tooltip} position={tooltipPosition} className="inline-flex">
			{button}
		</Tooltip>
	);
};

export default ActionButton;
