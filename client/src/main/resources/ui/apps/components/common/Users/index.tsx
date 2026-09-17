import { UserResponse } from '@/rest/userAPI';
import Minus from '@/components/@icons/minus';

export const userColumns = (t: (key: string) => string) => [
	{ title: t('label.login'), key: 'email', sortable: true },
	{
		title: t('label.organization'),
		key: 'name_organization',
		render: (u: UserResponse) => {
			return u.user_contact.name_organization !== '' ? (
				<span>{u.user_contact.name_organization}</span>
			) : (
				<div className="flex justify-center">
					<Minus fill="#ccc" size={20} />
				</div>
			);
		},
		sortable: true,
	},
	{
		title: t('label.address'),
		key: 'locality',
		render: (u: UserResponse) => {
			return u.user_contact.locality !== '' ? (
				<span>{u.user_contact.locality}</span>
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
		render: (u: UserResponse) => {
			const phone = u.user_contact.phone;
			return phone ? (
				<span>{`+${phone}`}</span>
			) : (
				<div className="flex justify-center">
					<Minus fill="#ccc" size={20} />
				</div>
			);
		},
		sortable: true,
	},
];
