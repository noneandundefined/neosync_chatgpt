import { useState } from 'react';
import PageLayout from '../PageLayout';
import { useTranslation } from 'react-i18next';
import { HELP_SECTIONs } from '@/constants/Help.constant';

import Help from '@/components/common/Help/Help';
import HelpUser from '@/components/common/Help/HelpUser';
import HelpGroup from '@/components/common/Help/HelpGroup';
import HelpDevice from '@/components/common/Help/HelpDevice';
import HelpCommand from '@/components/common/Help/HelpCommand';
import HelpSupport from '@/components/common/Help/HelpSupport';

const NeosyncHelpPage = () => {
	const { t } = useTranslation();

	const [section, setSection] = useState<string>('help');

	const getSectionText = (id: string) => {
		const helpSection = HELP_SECTIONs.find((s) => s.id === id);

		if (!helpSection) {
			return { title: '', description: '' };
		}

		return {
			title: t(`message.${helpSection.titleKey}`),
			description: t(`message.${helpSection.descriptionKey}`),
		};
	};

	return (
		<PageLayout>
			<div className="space-y-5 pb-3">
				{section === 'help' && <Help changeSection={setSection} />}

				{section === 'devices' && <HelpDevice {...getSectionText(section)} changeSection={() => setSection('help')} />}

				{section === 'commands' && <HelpCommand {...getSectionText(section)} changeSection={() => setSection('help')} />}

				{section === 'groups' && <HelpGroup {...getSectionText(section)} changeSection={() => setSection('help')} />}

				{section === 'users' && <HelpUser {...getSectionText(section)} changeSection={() => setSection('help')} />}

				<HelpSupport />
			</div>
		</PageLayout>
	);
};

export default NeosyncHelpPage;
