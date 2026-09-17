import IndexStep from './IndexStep';
import { useTranslation } from 'react-i18next';
import { UseFormRegister } from 'react-hook-form';
import GUIRadio from '@/components/ui/Radio/GUIRadio';
import GUISelect from '@/components/ui/Select/GUISelect';
import LabeledField from '@/components/ui/Form/LabeledField';
import { BULK_CONFIGURATION_MAX_DEVICES } from '@/constants/BulkConfiguration.constant';
import { CompanyCreateRequest } from '@/interface/company/companyCreateRequest.interface';
import { DeviceStatBC } from '@/pages/ConfigurationsPage/BulkConfigurationPage/BulkConfigurationPage';

interface ScopeOfApplicationStepProps {
	register: UseFormRegister<CompanyCreateRequest>;
	stats: DeviceStatBC;
	state: any;
	exceedsDeviceLimit?: boolean;
	onSelectModel: (model: string) => void;
}

const ScopeOfApplicationStep: React.FC<ScopeOfApplicationStepProps> = ({ register, stats, state, exceedsDeviceLimit = false, onSelectModel }) => {
	const { t } = useTranslation();

	return (
		<IndexStep step={2} title="message.provisioning-scope-of-application">
			<div className="border p-3 space-y-3 rounded-[6px]">
				<p className="font-medium text-[#49525f]">{t('message.provisioning-selected-devices-count', { count: state.selectedIds.length })}</p>

				{exceedsDeviceLimit && <p className="text-sm text-[#b42318]">{t('message.provisioning-company-max-devices-warning-short', { count: BULK_CONFIGURATION_MAX_DEVICES })}</p>}

				{stats.priorityModel && (
					<LabeledField label="message.provisioning-primary-model">
						<GUISelect value={stats.priorityModel} onChange={(event) => onSelectModel(event.target.value)}>
							{stats.models.map((model) => (
								<option key={model} value={model}>{t('message.provisioning-primary-model-option', { model, count: stats.modelCounts[model] })}</option>
							))}
						</GUISelect>
						{stats.dominantFirmware != null && stats.dominantFirmware > 0 && <p className="text-[13px] text-[#49525f]">{t('message.tracker-firmware-version')}: {`0x${stats.dominantFirmware.toString(16).toUpperCase().padStart(2, '0')}`}</p>}
					</LabeledField>
				)}

				{/* <div className="flex flex-wrap gap-3">
					{stats.models.map((model, index) => (
						<p key={index} className="text-sm text-[#49525f] border border-[#49525f] cursor-default px-3 py-1 rounded-[4px]">
							{model}
						</p>
					))}
				</div> */}
			</div>

			<LabeledField label="message.provisioning-incompatible-action-label">
				<GUIRadio label={t('message.provisioning-incompatible-always-excluded')} value="exclude" register={register('incompatible_action')} />
			</LabeledField>
		</IndexStep>
	);
};

export default ScopeOfApplicationStep;
