import { useEffect, useState } from 'react';
import { ConfigurationManager } from '@/utils/ConfigurationManagerUtils';

export const useConfigurationVal = (manager: ConfigurationManager | null, uid: string) => {
	const [value, setValue] = useState<any>(manager?.getField(uid)?.value);

	useEffect(() => {
		if (!manager) return;

		const update = () => {
			setValue(manager.getField(uid)?.value);
		};

		const unsubscribe = manager.subscribe(update);

		update();

		return () => {
			unsubscribe();
		};
	}, [manager, uid]);

	return value;
};
