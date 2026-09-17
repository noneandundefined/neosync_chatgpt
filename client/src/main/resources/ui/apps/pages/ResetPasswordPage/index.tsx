import React from 'react';
import ResetPasswordPage from './ResetPasswordPage';
import PageMeta from '@/components/PageMeta/PageMeta';

const IResetPasswordPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-reset-password" ogTitleKey="message.og-title-reset-password" ogDescriptionKey="message.og-description-reset-password" path="/password/new" />

			<ResetPasswordPage />
		</React.Fragment>
	);
};

export default IResetPasswordPage;
