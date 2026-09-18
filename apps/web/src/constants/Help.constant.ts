import Car from '@/components/@icons/car';
import Search from '@/components/@icons/search';
import Wrench from '@/components/@icons/wrench';
import Pencil from '@/components/@icons/pencil';
import Delete from '@/components/@icons/delete';
import TrainCar from '@/components/@icons/train-car';
import Information from '@/components/@icons/information';
import AccountGroup from '@/components/@icons/account-group';
import CloudOutline from '@/components/@icons/cloud-outline';
import AccountSwitch from '@/components/@icons/account-switch';
import CardBulletedOutline from '@/components/@icons/card-bulleted-outline';
import EmailArrowRight from '@/components/@icons/email-arrow-right';

export const HELP_SECTIONs = [
	{
		id: 'devices',
		titleKey: 'help-section-devices-title',
		descriptionKey: 'help-section-devices-description',
		icon: Car,
	},
	{
		id: 'commands',
		titleKey: 'help-section-commands-title',
		descriptionKey: 'help-section-commands-description',
		icon: CardBulletedOutline,
	},
	{
		id: 'groups',
		titleKey: 'help-section-groups-title',
		descriptionKey: 'help-section-groups-description',
		icon: TrainCar,
	},
	{
		id: 'users',
		titleKey: 'help-section-users-title',
		descriptionKey: 'help-section-users-description',
		icon: AccountGroup,
	},
];

export const HELP_DEVICEs_CONTROLs = [
	{
		icon: Search,
		titleKey: 'help-device-control-search-title',
		descriptionKey: 'help-device-control-search-description',
	},
	{
		icon: Information,
		titleKey: 'help-device-control-info-title',
		descriptionKey: 'help-device-control-info-description',
	},
	{
		icon: Wrench,
		titleKey: 'help-device-control-settings-title',
		descriptionKey: 'help-device-control-settings-description',
	},
	{
		icon: Pencil,
		titleKey: 'help-device-control-edit-title',
		descriptionKey: 'help-device-control-edit-description',
	},
	{
		icon: Delete,
		titleKey: 'help-device-control-delete-title',
		descriptionKey: 'help-device-control-delete-description',
	},
	{
		icon: CloudOutline,
		titleKey: 'help-device-control-connection-title',
		descriptionKey: 'help-device-control-connection-description',
	},
];

export const HELP_GROUPs_CONTROLs = [
	{
		icon: AccountSwitch,
		titleKey: 'help-group-control-copy-title',
		descriptionKey: 'help-group-control-copy-description',
	},
	{
		icon: Pencil,
		titleKey: 'help-group-control-edit-title',
		descriptionKey: 'help-group-control-edit-description',
	},
	{
		icon: Delete,
		titleKey: 'help-group-control-delete-title',
		descriptionKey: 'help-group-control-delete-description',
	},
];

export const HELP_COMMANDs_STEPs = [
	{
		id: 1,
		titleKey: 'help-command-step-immediate-title',
		descriptionKey: 'help-command-step-immediate-description',
	},
	{
		id: 2,
		titleKey: 'help-command-step-next-connection-title',
		descriptionKey: 'help-command-step-next-connection-description',
	},
];

export const HELP_COMMANDs_NOTE_KEY = 'help-command-note';
export const HELP_USERs_NOTE_KEY = 'help-user-note';

export const HELP_GROUPs_STEPs = [
	{
		id: 1,
		titleKey: 'help-group-step-create-title',
		descriptionKey: 'help-group-step-create-description',
	},
	{
		id: 2,
		titleKey: 'help-group-step-add-devices-title',
		descriptionKey: 'help-group-step-add-devices-description',
	},
	{
		id: 3,
		titleKey: 'help-group-step-send-commands-title',
		descriptionKey: 'help-group-step-send-commands-description',
	},
];

export const HELP_USERs_CONTROLs = [
	{
		icon: Information,
		titleKey: 'help-user-control-info-title',
		descriptionKey: 'help-user-control-info-desc',
	},
	{
		icon: EmailArrowRight,
		titleKey: 'help-user-control-email-title',
		descriptionKey: 'help-user-control-email-desc',
	},
	{
		icon: Pencil,
		titleKey: 'help-user-control-edit-title',
		descriptionKey: 'help-user-control-edit-desc',
	},
	{
		icon: Delete,
		titleKey: 'help-user-control-delete-title',
		descriptionKey: 'help-user-control-delete-desc',
	},
];

export const HELP_USERs_STEPs = [
	{
		id: 1,
		titleKey: 'help-user-step-create-title',
		descriptionKey: 'help-user-step-create-desc',
	},
	{
		id: 2,
		titleKey: 'help-user-step-manage-title',
		descriptionKey: 'help-user-step-manage-desc',
	},
	{
		id: 3,
		titleKey: 'help-user-step-filter-title',
		descriptionKey: 'help-user-step-filter-desc',
	},
];
