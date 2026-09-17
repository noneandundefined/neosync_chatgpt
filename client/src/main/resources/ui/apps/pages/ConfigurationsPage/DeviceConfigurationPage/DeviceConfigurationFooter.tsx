import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ModalConfirm from '@/components/Modal/ModalConfirm';
import ActionButton from '@/components/ui/Button/ActionButton';

import Broom from '@/components/@icons/broom';
import Reload from '@/components/@icons/reload';
import Factory from '@/components/@icons/factory';
import Download from '@/components/@icons/download';
import Save from '@/components/@icons/content-save';
import ChevronUp from '@/components/@icons/chevron-up';
import FolderOpen from '@/components/@icons/folder-open';
import ChevronDown from '@/components/@icons/chevron-down';
import ArrowCollapseUp from '@/components/@icons/arrow-collapse-up';

import Modal from '@/components/Modal/Modal';
import useIsMobile from '@/hooks/useIsMobile';
// import { formatDate } from '@/utils/TimeUtils';
import { useTcpHealth } from '@/context/useTcpHealth';
import { useDraftCheck } from '@/hooks/useDraftCheck';
import { basicUserGetMeAccesses } from '@/rest/userAPI';
// import Spinner from '@/components/common/Loader/Spinner';
// import GUISelect from '@/components/ui/Select/GUISelect';
import { useKeyboardOpen } from '@/hooks/useKeyboardOpen';
import { useModalContext } from '@/context/useModalContext';
import DatabaseSync from '@/components/@icons/database-sync';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { useConfigurationContext } from '@/context/useConfigurationContext';
// import { basicConfigurationGetHistories, basicConfigurationImportByHistory } from '@/rest/configurationAPI';

interface DeviceConfigurationFooterProps {
	imei: string;
	onImport: (file: File) => void;
	onExport: () => void;
	onApply: () => void;
	onReboot: () => void;
	onRefresh: () => void;
	onReloadSection: () => void;
	onEraseEeprom: () => void;
	onEraseFlash: () => void;
	onRebootTelemetry: () => void;
}

const DeviceConfigurationFooter: React.FC<DeviceConfigurationFooterProps> = ({ imei, onImport, onExport, onApply, onReboot, onRebootTelemetry, onRefresh, onEraseEeprom, onEraseFlash }) => {
	const { t } = useTranslation();

	const isMobile = useIsMobile();
	const tcpHealth = useTcpHealth();
	const keyboardOpen = useKeyboardOpen();

	const { device } = useConfigurationContext();

	const { open, close } = useModalContext();

	const { data: respUserGetMeAccesses } = useHandleServer(['respUserGetMeAccesses'], basicUserGetMeAccesses);

	const [panelOpen, setPanelOpen] = useState<boolean>(false);

	const isChanged = useDraftCheck(imei);
	const deviceOnline = tcpHealth && device.status;
	const deviceOfflineTooltip = !deviceOnline ? t('message.device-not-connected-to-server') : undefined;

	const [rebootLoading, setRebootLoading] = useState<boolean>(false);
	const [rebootTelemetryLoading, setRebootTelemetryLoading] = useState<boolean>(false);
	const [eraseEepromLoading, setEraseEepromLoading] = useState<boolean>(false);
	const [eraseFlashLoading, setEraseFlashLoading] = useState<boolean>(false);

	const [exportLoading, setExportLoading] = useState<boolean>(false);
	const [refreshLoading, setRefreshLoading] = useState<boolean>(false);
	const [applyLoading, setApplyLoading] = useState<boolean>(false);

	const openConfirmationModal = (title: string, message: string, onConfirm: () => void) => {
		open(
			<Modal title={title}>
				<ModalConfirm
					message={message}
					onConfirm={() => {
						onConfirm();
						close();
					}}
					onCancel={() => close()}
				/>
			</Modal>
		);
	};

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

	/** If keyboard open, close the panel */
	useEffect(() => {
		if (keyboardOpen && panelOpen) {
			setPanelOpen(false);
		}
	}, [keyboardOpen, panelOpen]);

	return (
		<>
			<div className="fixed bottom-0 left-0 z-[999] w-full">
				{isMobile && (
					<div id="buttonhlp" className={`flex justify-center !rounded-none cursor-pointer ${isChanged && !panelOpen && '!bg-[#7db500]'}`} onClick={() => setPanelOpen(!panelOpen)}>
						{panelOpen ? (
							<div className="flex items-end gap-1">
								<p className={isChanged && !panelOpen ? 'text-white' : ''}>{t('message.service-actions')}</p>
								<ChevronDown size={22} fill={isChanged && !panelOpen ? '#fff' : '#62666d'} />
							</div>
						) : (
							<div className="flex items-end gap-1">
								<p className={isChanged && !panelOpen ? 'text-white' : ''}>{t('message.service-actions')}</p>
								<ChevronUp size={22} fill={isChanged && !panelOpen ? '#fff' : '#62666d'} />
							</div>
						)}
					</div>
				)}

				{(panelOpen || !isMobile) && (
					<div className="px-2 sm:px-8 md:px-16 py-4 flex flex-col justify-between md:flex-row lg:items-start gap-[2rem] xl:gap-[4rem] bg-[#fafafa]">
						{(respUserGetMeAccesses?.access_configuration_apply ?? true) && (
							<div className="flex flex-row gap-1 sm:gap-6 justify-around lg:justify-start">
								<ActionButton icon={FolderOpen} label={t('message.read-settings-file')} disabled={!deviceOnline} tooltip={deviceOfflineTooltip} isChange={false} isFileInput={true} onFileSelect={onImport} />

								<ActionButton
									icon={Save}
									label={t('message.save-settings-file')}
									disabled={exportLoading}
									isChange={false}
									onClick={() => {
										openConfirmationModal(t('message.save-settings-file'), t('message.attention-neosync-type-file-export'), () => {
											setExportLoading(true);
											try {
												onExport();
												setPanelOpen(false);
											} finally {
												setExportLoading(false);
											}
										});
									}}
								/>

								<ActionButton
									icon={Download}
									label={t('message.read-settings')}
									disabled={!deviceOnline || refreshLoading}
									tooltip={deviceOfflineTooltip}
									isChange={false}
									onClick={() => {
										setRefreshLoading(true);
										try {
											onRefresh();
											setPanelOpen(false);
										} finally {
											setRefreshLoading(false);
										}
									}}
								/>

								<ActionButton
									icon={ArrowCollapseUp}
									label={t('message.apply-settings')}
									disabled={applyLoading}
									isChange={isChanged}
									onClick={() => {
										setApplyLoading(true);
										try {
											onApply();
											setPanelOpen(false);
										} finally {
											setApplyLoading(false);
										}
									}}
								/>
							</div>
						)}

						<div className="flex flex-row gap-10 lg:gap-6 justify-center lg:justify-end">
							<ActionButton
								icon={DatabaseSync}
								label={t('message.read-sensor-data')}
								disabled={!deviceOnline || rebootTelemetryLoading}
								tooltip={deviceOfflineTooltip}
								isChange={false}
								onClick={() => {
									setRebootTelemetryLoading(true);
									try {
										onRebootTelemetry();
									} finally {
										setRebootTelemetryLoading(false);
									}
								}}
							/>
							<ActionButton
								icon={Reload}
								label={t('label.reboot')}
								disabled={!deviceOnline || rebootLoading}
								tooltip={deviceOfflineTooltip}
								isChange={false}
								onClick={() =>
									openConfirmationModal(t('message.reboot-device-title'), t('message.reboot-device-confirmation'), () => {
										setRebootLoading(true);
										try {
											onReboot();
										} finally {
											setRebootLoading(false);
										}
									})
								}
							/>
							<ActionButton
								icon={Factory}
								label={t('label.factory-settings')}
								disabled={!deviceOnline || eraseEepromLoading}
								tooltip={deviceOfflineTooltip}
								isChange={false}
								onClick={() =>
									openConfirmationModal(t('message.factory-settings-title'), t('message.factory-settings-confirmation'), () => {
										setEraseEepromLoading(true);
										try {
											onEraseEeprom();
										} finally {
											setEraseEepromLoading(false);
										}
									})
								}
							/>
							<ActionButton
								icon={Broom}
								label={t('label.clear-memory')}
								disabled={!deviceOnline || eraseFlashLoading}
								tooltip={deviceOfflineTooltip}
								isChange={false}
								onClick={() =>
									openConfirmationModal(t('message.clear-memory-title'), t('message.clear-memory-confirmation'), () => {
										setEraseFlashLoading(true);
										try {
											onEraseFlash();
										} finally {
											setEraseFlashLoading(false);
										}
									})
								}
							/>
						</div>
					</div>
				)}
			</div>
		</>
	);
};

export default DeviceConfigurationFooter;
