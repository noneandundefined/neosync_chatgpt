import React from 'react';
import SignInPage from './SignInPage';
import PageMeta from '@/components/PageMeta/PageMeta';

const ISignInPage = () => {
	return (
		<React.Fragment>
			<PageMeta descriptionKey="message.meta-description-signin" ogTitleKey="message.og-title-signin" ogDescriptionKey="message.og-description-signin" path="/sign_in" />

			<SignInPage />
		</React.Fragment>
	);
};

export default ISignInPage;
