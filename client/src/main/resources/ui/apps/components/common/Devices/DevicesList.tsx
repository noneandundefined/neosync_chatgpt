import { useState } from 'react';
import { deviceColumns } from '.';
import { Link } from 'react-router-dom';
import { formatDate } from '@/utils/TimeUtils';
import { useTranslation } from 'react-i18next';
import EmptyDevices from './common/EmptyDevices';
import { UserAccessResponse } from '@/rest/userAPI';
import DeviceEditModal from './Modals/DeviceEditModal';
import React, { Dispatch, SetStateAction } from 'react';
import { buildRoute, ROUTES } from '@/constants/constants';
import { useModalContext } from '@/context/useModalContext';
import DeviceTransferModal from './Modals/DeviceTransferModal';
import ConfigurationSyncStatus from './common/ConfigurationSyncStatus';
import ModalConfirmRemoval from '@/components/Modal/ModalConfirmRemoval';
import GenericTable, { FetchFn } from '../Table/GenericTable/GenericTable';
import ModalUpdateFirmware from '@/components/Modal/ModalUpdateFirmware.tsx';
import { basicDeviceDelete, basicDeviceMassiveDelete, basicDevicesGetWithParams, Device_DeviceConf_Sync } from '@/rest/deviceAPI';

import Modal from '@/components/Modal/Modal';
import Wrench from '@/components/@icons/wrench';
import Delete from '@/components/@icons/delete';
import Pencil from '@/components/@icons/pencil';
import OpenInNew from '@/components/@icons/open-in-new';
import Information from '@/components/@icons/information';

export interface DevicesListProps {
	/** Профиль пользователя для доступов */
	userAccesses: UserAccessResponse | null;
	/** Токен для обновления данных в таблице */
	sysToken: number;
	/** Токен для обновления данных в таблице */
	setSysToken: Dispatch<SetStateAction<number>>;
}

const DevicesList: React.FC<DevicesListProps> = ({ userAccesses, sysToken, setSysToken }) => {
	const { t } = useTranslation();

	const { open, close } = useModalContext();

	const [expandedRow, setExpandedRow] = useState<number | null>(null);
	const [updatedDevices, setUpdatedDevices] = useState<Set<string>>(new Set());

	const handleDelete = async (payloads: (string | number)[]) => {
		await basicDeviceMassiveDelete(payloads as string[]);
		setSysToken((v: number) => v + 1);
	};

	const handleRowAction = (device: Device_DeviceConf_Sync, idx: number, isHovered: boolean) => (
		<div
			className="flex items-center justify-center gap-2 sm:gap-1"
			style={{
				opacity: isHovered ? 1 : 0,
			}}
		>
			<Information
				fill="#1976d2"
				size={24}
				className="cursor-pointer transition"
				onClick={(e?: React.MouseEvent<SVGSVGElement>) => {
					e?.stopPropagation();
					setExpandedRow(expandedRow === idx ? null : idx);
				}}
			/>

			{(userAccesses?.access_configuration_read ?? true) && (
				<Link
					to={buildRoute(ROUTES.CONFIGURATIONS_DEVICE, {
						imei: device.imei,
					})}
					className="relative"
				>
					<div>
						<Wrench fill="#ff7c00" size={24} className="cursor-pointer transition" />
					</div>
				</Link>
			)}

			{(userAccesses?.access_treker_delete ?? true) && (
				<React.Fragment>
					<Pencil
						fill="#5ba0e8"
						size={24}
						className="cursor-pointer transition"
						onClick={() => {
							open(
								<Modal title={t('message.edit-device')}>
									<DeviceEditModal
										imei={device.imei}
										onSuccess={() => {
											setSysToken((v: number) => v + 1);
											close();
										}}
									/>
								</Modal>
							);
						}}
					/>

					<Delete
						fill="#f40000"
						size={24}
						className="cursor-pointer transition"
						onClick={() => {
							open(
								<Modal title={t('message.delete-title-device-sure')}>
									<ModalConfirmRemoval
										object={device.imei}
										onDelete={async () => {
											await basicDeviceDelete(device.imei);
											setSysToken((v) => v + 1);
											close();
										}}
									/>
								</Modal>
							);
						}}
					/>
				</React.Fragment>
			)}
		</div>
	);

	const renderExpandedRow = (device: Device_DeviceConf_Sync, idx: number) =>
		expandedRow === idx ? (
			<tr className="renderExpandedRow">
				<td colSpan={11} className="!w-full">
					<div className="p-1 text-left flex items-center justify-between w-full">
						<ul className="space-y-2">
							{(userAccesses?.access_log_read ?? true) && <li>{t('message.detailed-device-logs')}</li>}
							<li>{t('message.remote-configuration-server')}</li>
							{device.activated && (
								<>
									<li>{t('message.configuration-sync-status-title')}</li>
									<li>{t('message.date-last-connection')}</li>
								</>
							)}
							{device.activated && (
								<>
									<li>{t('message.last-setup-date')}</li>
									<li>{t('message.device-update-date')}</li>
								</>
							)}
						</ul>

						<ul className="space-y-2 text-right">
							{(userAccesses?.access_log_read ?? true) && (
								<li>
									<Link
										to={buildRoute(ROUTES.DEVICE_LOGS, {
											imei: device.imei,
										})}
										className="flex justify-end items-center gap-1 text-[#000] hover:text-[#000]"
									>
										<p className="font-normal">{t('label.open')}</p>
										<OpenInNew fill="#000" size={15} />
									</Link>
								</li>
							)}
							<li className={`${device.status ? 'text-[#3a9d08]' : device.activated ? 'text-[#9d0808]' : 'text-[#dd6300]'}`}>
								{device.status ? t('message.packet-has-been-sent-successfully') : device.activated ? t('label.disabled') : t('message.connection-expected')}
							</li>
							{device.activated && (
								<>
									<li>
										<ConfigurationSyncStatus device={device} />
									</li>
									<li>{device.last_mod_time !== 0 ? formatDate(device.updated_at) : t('message.no-information-available')}</li>
								</>
							)}
							{device.activated && (
								<>
									<li>{device.configuration_updated_at ? formatDate(device.configuration_updated_at) : t('message.no-information-available')}</li>
									<li>
										{device.firmware_update_status === 'pending'
											? t('message.firmware-update-in-progress')
											: device.firmware_updated_at
												? formatDate(device.firmware_updated_at)
												: t('message.no-information-available')}
									</li>
								</>
							)}
						</ul>
					</div>
				</td>
			</tr>
		) : null;

	const columns = deviceColumns(
		(device) => {
			open(
				<Modal title={t('message.update-firmware')}>
					<ModalUpdateFirmware device={device} onSuccess={(imei) => setUpdatedDevices((prev) => new Set(prev).add(imei))} />
				</Modal>
			);
		},
		updatedDevices,
		setSysToken
	);

	const fetchDevices: FetchFn<Device_DeviceConf_Sync> = async ({ page, limit, search, columnSort }, signal) => {
		const res = await basicDevicesGetWithParams(page, limit, search, columnSort, signal);
		return { items: res.items, total: res.total, totalPages: res.total_pages };
	};

	const handleTransferSelected = async (ids: (string | number)[]) => {
		const imeis = ids.map(String).filter(Boolean);
		if (imeis.length === 0) {
			return;
		}

		open(
			<Modal title={t('label.transfer-title-device-sure')}>
				<DeviceTransferModal
					imeis={imeis}
					onSuccess={() => {
						setSysToken((v: number) => v + 1);
						close();
					}}
				/>
			</Modal>
		);
	};

	return (
		<React.Fragment>
			<GenericTable<Device_DeviceConf_Sync>
				tableKey="fetchDevices"
				columns={columns}
				fetchFn={fetchDevices}
				sysToken={sysToken}
				getRowId={(device) => device.imei}
				rowActions={(device, idx, isHovered) => handleRowAction(device, idx, isHovered)}
				renderExpandedRow={renderExpandedRow}
				onDeleteFn={handleDelete}
				enableTransferSelectedAction={true}
				onTransferSelected={handleTransferSelected}
				emptyComponent={<EmptyDevices userAccesses={userAccesses} setSysToken={setSysToken} />}
			/>
		</React.Fragment>
	);
};

export default DevicesList;
