import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import ADM40 from './ADM40';
import Period from './Period';
import NavTimeSync from './NavTimeSync';
import NavigationFilter from './NavigationFilter';
import TrackSavingOptions from './TrackSavingOptions';
import FilteringEmissions from './FilteringEmissions';
import IndexConfiguration from '../IndexConfiguration';
import FreezingCoordinates from './FreezingCoordinates';
import AlternativeSourcesPositioning from './AlternativeSourcesPositioning';
import SubstitutionCoordinatesAltSource from './SubstitutionCoordinatesAltSource';
import { DEVICE_MODELS, getStaticSupport } from '@/constants/DeviceModels.constant';
import { CACHEKEYs_TERMINAL_TRACK_EXPERT_SETTINGs } from '@/constants/CacheKeys.constants';
import { ConfigurationWrapperProvider, useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const ConfigTrackInner = () => {
	const { t } = useTranslation();

	const { imei, support, model, isTemplate } = useConfigurationWrapperContext();
	const staticSupport = getStaticSupport((model ?? DEVICE_MODELS.ADM333V2) as keyof typeof DEVICE_MODELS, isTemplate);

	const [expertSettings, setExpertSettings] = useState<boolean>(() => {
		const stored = sessionStorage.getItem(CACHEKEYs_TERMINAL_TRACK_EXPERT_SETTINGs(imei));
		return stored === 'true';
	});

	useEffect(() => {
		sessionStorage.setItem(CACHEKEYs_TERMINAL_TRACK_EXPERT_SETTINGs(imei), String(expertSettings));
	}, [expertSettings]);

	return (
		<div className="flex flex-col">
			<div className="flex flex-col sm:flex-row flex-wrap gap-x-4 items-stretch">
				<IndexConfiguration title={t('label.period')} content={<Period />} />
				<IndexConfiguration title={t('message.coordinate-freeze-at-standstill')} content={<FreezingCoordinates />} />
			</div>

			<IndexConfiguration title={t('message.track-save-parameters')} content={<TrackSavingOptions expertSettings={expertSettings} setExpertSettings={setExpertSettings} />} />

			<div className="flex flex-col sm:flex-row flex-wrap gap-x-4 items-stretch">
				{support.NavigationFilter && <IndexConfiguration title={t('message.configuring-navigation-filter')} content={<NavigationFilter />} />}

				{support.AlternativeSourcesPositioning && (isTemplate || staticSupport.AlternativeSourcesPositioning) && (
					<IndexConfiguration title={t('message.alternative-sources-positioning')} content={<AlternativeSourcesPositioning />} />
				)}

				{support.Navtimesync && <IndexConfiguration title={t('message.navtimesync')} content={<NavTimeSync />} />}

				{support.ADM40 && (isTemplate || !support.ModeHybrid) && <IndexConfiguration title="ADM40" content={<ADM40 />} />}

				{support.FilteringEmissions && <IndexConfiguration title={t('message.filtering-emissions')} content={<FilteringEmissions />} />}

				{support.SubstitutionCoordinatesAlternativeSource && <IndexConfiguration title={t('message.substitution-coordinates-alternative-source')} content={<SubstitutionCoordinatesAltSource />} />}
			</div>
		</div>
	);
};

const ConfigTrack = () => {
	return (
		<ConfigurationWrapperProvider>
			<ConfigTrackInner />
		</ConfigurationWrapperProvider>
	);
};

export default ConfigTrack;
