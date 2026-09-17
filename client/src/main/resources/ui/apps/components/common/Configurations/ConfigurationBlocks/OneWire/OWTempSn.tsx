import React from 'react';
import { useEffect, useState } from 'react';
import Close from '@/components/@icons/close';
import { useTranslation } from 'react-i18next';
import { UIDs } from '@/constants/UID.constant';
import Colors from '@/constants/Color.constant';
import { useConfigurationField } from '@/hooks/Configuration/useConfigurationField';
import { useConfigurationWrapperContext } from '@/context/useConfigurationWrapperContext';

const OWTempSn = () => {
	const { t } = useTranslation();
	const { imei, section, manager } = useConfigurationWrapperContext();

	const { value: owTempSn, handleChange: handleOwTempSn, saveDraft: saveOwTempSn } = useConfigurationField(imei, section, manager, UIDs.OW_TEMP_SN);

	const [owTemps, setOwTemps] = useState<string[] | null>(owTempSn);

	useEffect(() => {
		if (owTempSn) {
			setOwTemps(owTempSn);
		}
	}, [owTempSn]);

	const handleRemove = async (index: number) => {
		setOwTemps((prev) => {
			if (!prev) return null;

			const updatedOwTemps = prev.filter((_, i) => i !== index);
			handleOwTempSn(updatedOwTemps);
			saveOwTempSn(updatedOwTemps);
			return updatedOwTemps;
		});
	};

	return (
		<div className="space-y-5">
			{owTemps && owTemps.length > 0 && (
				<>
					{owTemps.map((owTemp, index) => {
						return (
							<div key={index}>
								<React.Fragment>
									<div className="border border-[#e5e8eb] p-1 px-2">
										<div className="flex justify-between items-center min-w-[17rem]">
											<span className="text-[13px]">
												{index} - {'Undefined'}
											</span>
											<div className="flex items-center gap-3">
												<span className="text-[13px]">{owTemp}</span>
												<span onClick={() => handleRemove(index)}>
													<Close size={22} fill={Colors.color_error_base} className="p-[1px] border border-[#dedede] cursor-pointer hover:bg-[#eee] transition" onClick={() => {}} />
												</span>
											</div>
										</div>
									</div>
								</React.Fragment>
							</div>
						);
					})}

					<div className="min-w-[10rem]">
						<button className="opacity-70 cursor-not-allowed bg-[#eee] text-[10px] text-black font-medium p-1 rounded-[6px] uppercase text-center hover:bg-[#e3e3e3] transition border border-[#ccc]">
							{t('message.activate-data-request')}
						</button>
					</div>
				</>
			)}
		</div>
	);
};

export default OWTempSn;
