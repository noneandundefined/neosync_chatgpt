import { useEffect, useState } from 'react';
import Modal from '@/components/Modal/Modal';
import useIsMobile from '@/hooks/useIsMobile';
import { useTranslation } from 'react-i18next';
import { useQueryClient } from '@tanstack/react-query';
import { useKeyboardOpen } from '@/hooks/useKeyboardOpen';
import { useModalContext } from '@/context/useModalContext';
import ActionButton from '@/components/ui/Button/ActionButton';
import { saveExistingTemplateDraft } from '@/utils/saveTemplateConfiguration';
import { useTemplateConfigurationContext } from '@/context/useTemplateConfigurationContext';
import ModalSaveTemplate from '@/components/common/Configurations/ConfigurationTemplate/Modals/ModalSaveTemplate';
import ModalRemoveTemplate from '@/components/common/Configurations/ConfigurationTemplate/Modals/ModalRemoveTemplate';

import ChevronUp from '@/components/@icons/chevron-up';
import ChevronDown from '@/components/@icons/chevron-down';
import OfficeBuildingCogOutline from '@/components/@icons/office-building-cog-outline';
import OfficeBuildingRemoveOutline from '@/components/@icons/office-building-remove-outline';

interface ConfigurationTemplateFooterProps {}

const ConfigurationTemplateFooter: React.FC<ConfigurationTemplateFooterProps> = () => {
	const { t } = useTranslation();

	const isMobile = useIsMobile();

	const keyboardOpen = useKeyboardOpen();

	const queryClient = useQueryClient();
	const templateCtx = useTemplateConfigurationContext();

	const isEditMode = !!templateCtx?.templateId;

	const [panelOpen, setPanelOpen] = useState<boolean>(false);

	const { open, close } = useModalContext();

	const [saveLoading, setSaveLoading] = useState<boolean>(false);

	useEffect(() => {
		const padding = isMobile ? (panelOpen ? 18 : 3) : 10;
		const paddingPx = padding * 16;

		document.documentElement.style.setProperty('--bottom-panel-padding', `${paddingPx}px`);

		const offset = 230;
		const duration = 300;
		const startScroll = window.scrollY;
		const targetScroll = panelOpen ? startScroll + offset : startScroll - offset;
		const startTime = performance.now();

		const animate = (time: number) => {
			const elapsed = time - startTime;
			const progress = Math.min(elapsed / duration, 1);

			const eased = progress < 0.5 ? 2 * progress * progress : -1 + (4 - 2 * progress) * progress;

			window.scrollTo(0, startScroll + (targetScroll - startScroll) * eased);

			if (progress < 1) requestAnimationFrame(animate);
		};

		requestAnimationFrame(animate);
	}, [panelOpen, keyboardOpen]);

	useEffect(() => {
		if (keyboardOpen && panelOpen) {
			setPanelOpen(false);
		}
	}, [keyboardOpen, panelOpen]);

	const handleSave = async () => {
		if (isEditMode && templateCtx?.templateId && templateCtx.template) {
			try {
				setSaveLoading(true);
				await saveExistingTemplateDraft(templateCtx.templateId, templateCtx.template, queryClient);
			} finally {
				setSaveLoading(false);
			}

			return;
		}

		open(
			<Modal title={t('message.saving-configuration-template')}>
				<ModalSaveTemplate saveLoading={saveLoading} setSaveLoading={setSaveLoading} onClose={close} />
			</Modal>
		);
	};

	return (
		<>
			<div className="fixed bottom-0 left-0 z-[999] w-full">
				{isMobile && (
					<div id="buttonhlp" className={`flex justify-center !rounded-none cursor-pointer`} onClick={() => setPanelOpen(!panelOpen)}>
						{panelOpen ? (
							<div className="flex items-end gap-1">
								<p className="text-[#62666d]">{t('message.service-actions')}</p>
								<ChevronDown size={22} fill="#62666d" />
							</div>
						) : (
							<div className="flex items-end gap-1">
								<p className="text-[#62666d]">{t('message.service-actions')}</p>
								<ChevronUp size={22} fill="#62666d" />
							</div>
						)}
					</div>
				)}

				{(panelOpen || !isMobile) && (
					<div className="px-2 sm:px-8 md:px-16 py-4 flex flex-col justify-between sm:flex-row lg:items-start gap-[2rem] xl:gap-[4rem] bg-[#fafafa]">
						<div className="flex flex-row gap-1 sm:gap-6 justify-around lg:justify-start">
							<ActionButton icon={OfficeBuildingCogOutline} label={isEditMode ? t('label.save') : t('label.save-configuration-template')} disabled={saveLoading} isChange={false} onClick={handleSave} />

							<ActionButton
								icon={OfficeBuildingRemoveOutline}
								label={t('label.remove-configuration-template')}
								disabled={saveLoading}
								isChange={false}
								onClick={() => {
									open(
										<Modal title={t('message.remove-configuration-template-title')}>
											<ModalRemoveTemplate saveLoading={saveLoading} setSaveLoading={setSaveLoading} onClose={close} />
										</Modal>
									);
								}}
							/>
						</div>
					</div>
				)}
			</div>
		</>
	);
};

export default ConfigurationTemplateFooter;
