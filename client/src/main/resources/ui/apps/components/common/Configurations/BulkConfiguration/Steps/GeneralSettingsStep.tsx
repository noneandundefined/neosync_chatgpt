import IndexStep from './IndexStep';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import AutoMode from '@/components/@icons/auto-mode';
import { GUInput } from '@/components/ui/Input/GUInput';
import GUISelect from '@/components/ui/Select/GUISelect';
import LabeledField from '@/components/ui/Form/LabeledField';
import { Control, Controller, FieldErrors, UseFormRegister } from 'react-hook-form';
import { CompanyCreateRequest } from '@/interface/company/companyCreateRequest.interface';

interface GeneralSettingsStepProps {
	register: UseFormRegister<CompanyCreateRequest>;
	control: Control<CompanyCreateRequest>;
	errors: FieldErrors<CompanyCreateRequest>;
	onGenerateName: () => void;
}

const GeneralSettingsStep: React.FC<GeneralSettingsStepProps> = ({ register, control, errors, onGenerateName }) => {
	const { t } = useTranslation();

	return (
		<IndexStep step={1} title="label.general">
			<LabeledField label="label.title" argv="*">
				<div className="flex items-start gap-3">
					<div className="min-w-0 flex-1">
						<GUInput
							type="text"
							{...register('name', {
								required: t('message.validation-required-field'),
							})}
							error={errors.name?.message}
						/>
					</div>
					<Tooltip title={t('message.provisioning-generate-name')} position="left" className="shrink-0">
						<button type="button" className="p-2 hover:bg-white rounded-full" onClick={onGenerateName} aria-label={t('message.provisioning-generate-name')}>
							<AutoMode fill="#49525f" size={22} />
						</button>
					</Tooltip>
				</div>
			</LabeledField>

			<LabeledField label="message.provisioning-task-ttl">
				<Controller
					name="ttl"
					control={control}
					render={({ field }) => (
						<GUISelect value={field.value} onChange={field.onChange}>
							<option value="1">{t('message.provisioning-ttl-1-day')}</option>
							<option value="7">{t('message.provisioning-ttl-7-days')}</option>
							<option value="30">{t('message.provisioning-ttl-30-days')}</option>
						</GUISelect>
					)}
				/>
			</LabeledField>

			<LabeledField label="message.provisioning-existing-task-action">
				<Controller
					name="existing_task_action"
					control={control}
					render={({ field }) => (
						<GUISelect value={field.value} onChange={field.onChange}>
							<option value="skip">{t('message.provisioning-action-skip')}</option>
							<option value="replace">{t('message.provisioning-action-replace')}</option>
						</GUISelect>
					)}
				/>
			</LabeledField>

			<LabeledField label="message.provisioning-launch-mode">
				<Controller
					name="launch_mode"
					control={control}
					render={({ field }) => (
						<GUISelect value={field.value} onChange={field.onChange}>
							<option value="all">{t('message.provisioning-launch-mode-all')}</option>
							{field.value === 'selected' && <option value="selected">{t('message.provisioning-launch-mode-selected')}</option>}
						</GUISelect>
					)}
				/>
			</LabeledField>
		</IndexStep>
	);
};

export default GeneralSettingsStep;
