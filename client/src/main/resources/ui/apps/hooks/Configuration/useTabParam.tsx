import { StepKey } from '@/constants/StepKeys.constant';
import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';

export const useTabParam = (defaultKey: StepKey = 'device') => {
	const [searchParams, setSearchParams] = useSearchParams();

	const tabParam = (searchParams.get('tab') as StepKey | null) ?? defaultKey;
	const [currentStep, setCurrentStepState] = useState<StepKey>(tabParam);
	const [manualChange, setManualChange] = useState(false);

	const clearTabParam = () => {
		searchParams.delete('tab');
		setSearchParams(searchParams, { replace: true });
	};

	const setCurrentStep = (step: StepKey, manual: boolean = true) => {
		setManualChange(manual);
		setCurrentStepState(step);
	};

	useEffect(() => {
		if (!manualChange) return;
		setSearchParams({ tab: currentStep }, { replace: true });
		setManualChange(false);
	}, [currentStep, manualChange, setSearchParams]);

	useEffect(() => {
		if (!manualChange) {
			setCurrentStepState(tabParam);
		}
	}, [tabParam, manualChange]);

	return { currentStep, setCurrentStep, clearTabParam };
};
