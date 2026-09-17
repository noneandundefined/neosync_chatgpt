import { useTranslation } from 'react-i18next';

interface GUISwitchProps extends React.InputHTMLAttributes<HTMLInputElement> {}

const GUISwitch: React.FC<GUISwitchProps> = ({ className, ...props }) => {
	const { t } = useTranslation();

	return (
		<div className="flex gap-2 items-center">
			<p className="w-9 text-right text-sm text-[#6b747e] cursor-default">{props.checked ? t('label.on') : t('label.off')}</p>
			<input type="checkbox" role="switch" className={`check-custom ${className || ''}`} {...props} />
		</div>
	);
};

export default GUISwitch;
