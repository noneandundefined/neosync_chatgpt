import PageLayout from '@/pages/PageLayout';
import { Helmet } from 'react-helmet-async';
import { useCallback, useState } from 'react';
import Tooltip from '@/components/ui/Tooltip';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import Fallback from '@/components/Fallback/Fallback';
import { useNavigate, useParams } from 'react-router-dom';
import { Loading } from '@/components/common/Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';
import { basicDeviceCommandById } from '@/rest/deviceCommandAPI';
import { formatCommandResponse } from '@/utils/CommandResponseUtils';

import Minus from '@/components/@icons/minus';
import Check from '@/components/@icons/check';
import Close from '@/components/@icons/close';
import ArrowLeft from '@/components/@icons/arrow-left';
import ChevronDown from '@/components/@icons/chevron-down';
import ChevronRight from '@/components/@icons/chevron-right';

const DeviceCommandDetailsPage = () => {
	const { t } = useTranslation();
	const navigate = useNavigate();

	const { id } = useParams();
	if (!id) {
		return <Fallback />;
	}

	const [openIndex, setOpenIndex] = useState<number | null>(null);

	const fetchCommandById = useCallback(() => basicDeviceCommandById(Number(id)), [Number(id)]);
	const { data: respDeviceCommandById, loading: respDeviceCommandByIdLoading } = useHandleServer(['respDeviceCommandById', Number(id)], fetchCommandById, { refetchInterval: 3000 });

	if (!respDeviceCommandById) navigate(ROUTES.NOT_FOUND);

	return (
		<PageLayout>
			{respDeviceCommandById && (
				<Helmet>
					<title>{`${t('label.page-details-command')} · ${respDeviceCommandById.command} #${respDeviceCommandById.id}`}</title>
				</Helmet>
			)}

			<div className="flex-1 bg-white border border-[#e4e4e4] py-4 space-y-5 rounded flex flex-col overflow-y-auto">
				<div className="flex items-center px-5 md:px-9">
					<div className="flex flex-col space-y-2">
						<div className="flex items-center gap-2 -ml-2 p-1 px-2 rounded-md self-start hover:bg-[#f3f3f3] cursor-pointer" onClick={() => navigate(-1)}>
							<ArrowLeft fill="#999" size={15} />
							<p className="text-[#999] text-[13px]">{t('label.commands')}</p>
						</div>

						<div className="flex items-center gap-3">
							<div
								className="rounded-full p-[3px]"
								style={{
									background:
										respDeviceCommandById?.status === 'executionerror'
											? '#e60c00'
											: respDeviceCommandById?.status === 'completed'
												? '#029b00'
												: respDeviceCommandById?.status === 'notcompleted'
													? '#49525f'
													: 'transparent',
								}}
							>
								{respDeviceCommandById?.status === 'completed' && (
									<Tooltip title={t('label.status-success')}>
										<Check size={15} fill="#fff" />
									</Tooltip>
								)}
								{respDeviceCommandById?.status === 'notcompleted' && (
									<Tooltip title={t('label.notcompleted')}>
										<div className="rotate-45">
											<Minus size={15} fill="#ffffff" />
										</div>
									</Tooltip>
								)}
								{respDeviceCommandById?.status === 'executionerror' && (
									<Tooltip title={t('label.status-failed')}>
										<Close size={15} fill="#fff" />
									</Tooltip>
								)}
								{respDeviceCommandById?.status === 'inprogress' && (
									<Tooltip title={t('label.status-in-progress')}>
										<div className="relative flex justify-center items-center">
											<div className="w-[12px] h-[12px] rounded-full bg-[#9a6700]" />
											<div className={`absolute w-[12px] h-[12px] border-b-2 border-[#9a6700] rounded-full p-[9px] animate-spin`}></div>
										</div>
									</Tooltip>
								)}
								{respDeviceCommandById?.status === 'pending' && (
									<Tooltip title={t('label.status-pending')}>
										<div className="relative flex justify-center items-center">
											<div className="w-[12px] h-[12px] rounded-full bg-[#9a6700]" />
										</div>
									</Tooltip>
								)}
							</div>

							<p className="font-medium text-[17px] uppercase">
								{respDeviceCommandById?.command}
								<span className="font-normal text-[#999] ml-2">#{respDeviceCommandById?.id}</span>
							</p>
						</div>
					</div>
				</div>

				<div className="w-full h-[1px] bg-[#e4e4e4]" />

				<div className="space-y-3">
					{respDeviceCommandByIdLoading ? (
						<Loading />
					) : (
						respDeviceCommandById?.executions.map((device, index) => {
							const isOpen = openIndex === index;

							return (
								<div key={index} className="px-3 md:px-6">
									<div className="flex items-center gap-3 hover:bg-[#f3f3f3] p-2 rounded-lg cursor-pointer" onClick={() => setOpenIndex(isOpen ? null : index)}>
										{isOpen ? <ChevronDown fill="#49525f" size={22} /> : <ChevronRight fill="#49525f" size={22} />}

										<div className="flex items-center gap-3">
											<div>
												{device.status === 'completed' && (
													<Tooltip title={t('label.status-success')}>
														<div className="rounded-full p-[1px] bg-[#49525f]">
															<Check size={13} fill="#fff" />
														</div>
													</Tooltip>
												)}
												{device.status === 'notcompleted' && (
													<Tooltip title={t('label.notcompleted')}>
														<div className="rotate-45 rounded-full p-[1px] bg-[#49525f]">
															<Minus size={15} fill="#fff" />
														</div>
													</Tooltip>
												)}
												{device.status === 'executionerror' && (
													<Tooltip title={t('label.status-failed')}>
														<div className="rounded-full p-[1px] bg-[#49525f]">
															<Close size={13} fill="#fff" />
														</div>
													</Tooltip>
												)}
												{device.status === 'inprogress' && (
													<Tooltip title={t('label.status-in-progress')}>
														<div className="relative flex justify-center items-center">
															<div className="w-[10px] h-[10px] rounded-full bg-[#9a6700]" />
															<div className={`absolute w-[12px] h-[12px] border-b-2 border-[#9a6700] rounded-full p-[7px] animate-spin`}></div>
														</div>
													</Tooltip>
												)}
												{device.status === 'pending' && (
													<Tooltip title={t('label.status-pending')}>
														<div className="relative flex justify-center items-center">
															<div className="w-[10px] h-[10px] rounded-full bg-[#9a6700]" />
														</div>
													</Tooltip>
												)}
											</div>

											<p className="text-[15px] text-[#49525f]">{device.imei}</p>
										</div>
									</div>

									{isOpen && (
										<div className="hover:bg-[#f3f3f3] px-[2rem] md:px-[4rem] my-3">
											<p className="text-[13px]">{formatCommandResponse(device.response, t)}</p>
										</div>
									)}
								</div>
							);
						})
					)}
				</div>
			</div>
		</PageLayout>
	);
};

export default DeviceCommandDetailsPage;
