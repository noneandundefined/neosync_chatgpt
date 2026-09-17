import { UserResponse } from '@/rest/userAPI';
import { formatDate } from '@/utils/TimeUtils';

interface UserExpandedRowProps {
	user: UserResponse;
	t: any;
}

const UserExpandedRow: React.FC<UserExpandedRowProps> = ({ user, t }) => {
	return (
		<tr className="renderExpandedRow">
			<td colSpan={11} className="!w-full">
				<div className="p-1 text-left flex items-center justify-between w-full">
					<ul className="space-y-2">
						<li>{t('message.created-admin')}</li>
						<li>{t('label.created')}</li>
					</ul>

					<ul className="space-y-2 text-right">
						<li>{user.parent_email ? user.parent_email : t('label.unknown')}</li>
						<li>{formatDate(user.user_contact.created_at)}</li>
					</ul>
				</div>
			</td>
		</tr>
	);
};

export default UserExpandedRow;
