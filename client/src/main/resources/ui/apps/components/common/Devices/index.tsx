import Owner from './common/Owner';
import Minus from '@/components/@icons/minus';
import { useTranslation } from 'react-i18next';
import FirmwareCell from './common/FirmwareCell';
import { Dispatch, SetStateAction } from 'react';
import { Device_DeviceConf_Sync } from '@/rest/deviceAPI';
import { hasPermission } from '@/constants/Roles.constant';
import CloudOutline from '@/components/@icons/cloud-outline';
import { useRole } from '@/context/RoleContext/useRoleContext';
import CloudOffOutline from '@/components/@icons/cloud-off-outline';
import { convertUnixTimestampToPrettyFormat } from '@/utils/UnixConvertUtils';

export const deviceColumns = (onFirmwareClick: (device: Device_DeviceConf_Sync) => void, updatedDevices: Set<string>, setSysToken: Dispatch<SetStateAction<number>>) => {
	const { role } = useRole();
	const { t } = useTranslation();

	return [
		{ title: t('label.imei'), key: 'imei', sortable: true },
		{
			title: t('label.model'),
			key: 'device_model',
			render: (d: Device_DeviceConf_Sync) => {
				return <span className={d.device_model ? `text-black` : `text-[#ccc]`}>{d.device_model ? d.device_model : t('label.unknown')}</span>;
			},
			sortable: true,
		},
		{
			title: t('label.firmware'),
			key: 'firmware_version',
			render: (d: Device_DeviceConf_Sync) => <FirmwareCell d={d} t={t} onClick={onFirmwareClick} updatedDevices={updatedDevices} />,
			sortable: true,
		},
		// { title: t('label.groups'), key: 'group_name' },
		{
			title: t('label.name'),
			key: 'name',
			render: (d: Device_DeviceConf_Sync) => {
				return d.name ? (
					<span>{d.name}</span>
				) : (
					<div className="flex justify-center">
						<Minus fill="#ccc" size={20} />
					</div>
				);
			},
			sortable: true,
		},
		{
			title: t('label.phone'),
			key: 'phone',
			render: (d: Device_DeviceConf_Sync) => {
				return d.phone ? (
					<span>{`+${d.phone}`}</span>
				) : (
					<div className="flex justify-center">
						<Minus fill="#ccc" size={20} />
					</div>
				);
			},
			sortable: true,
		},
		{
			title: t('label.organization'),
			key: 'name_organization',
			render: (d: Device_DeviceConf_Sync) => {
				return d.name_organization ? (
					<span>{d.name_organization}</span>
				) : (
					<div className="flex justify-center">
						<Minus fill="#ccc" size={20} />
					</div>
				);
			},
			sortable: true,
		},
		{
			title: t('label.date-sync'),
			key: 'last_mod_time',
			render: (d: Device_DeviceConf_Sync) => {
				const formatted = convertUnixTimestampToPrettyFormat(d.last_mod_time);
				const [date, time] = formatted.split(' ');

				return d.last_mod_time !== 0 ? (
					<span className="flex flex-col max-[1250px]:flex-row max-[1250px]:gap-1">
						<span>{date}</span>
						<span>{time}</span>
					</span>
				) : (
					<span className="text-[#ccc]">{t('label.unknown')}</span>
				);
			},
			sortable: true,
		},
		// Используем оператор распространения (...), чтобы добавить элемент массива только при выполнении условия: hashPermission
		...(role !== null && hasPermission(role, 'transfer:users')
			? [
				{
					title: t('label.owner'),
					key: 'email',
					render: (d: Device_DeviceConf_Sync) => <Owner d={d} onUpdate={() => setSysToken((v) => v + 1)} />,
					sortable: true,
				},
			]
			: []),
		{
			title: t('label.connection'),
			key: 'status',
			render: (d: Device_DeviceConf_Sync) =>
				d.status ? (
					<span className="flex justify-center">
						<CloudOutline fill="#43c800" size={24} />
					</span>
				) : (
					<span className="flex justify-center">
						<CloudOffOutline fill="#aaa" size={24} />
					</span>
				),
			sortable: false,
		},
	];
};
