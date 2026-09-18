import React from 'react';
import PageMeta from '@/components/PageMeta/PageMeta';
import ReqResetPasswordPage from './ReqResetPasswordPage';

const IReqResetPasswordPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-reset-password" ogTitleKey="message.og-title-reset-password" ogDescriptionKey="message.og-description-reset-password" path="/password_reset" />

			<ReqResetPasswordPage />
		</React.Fragment>
	);
};

export default IReqResetPasswordPage;
