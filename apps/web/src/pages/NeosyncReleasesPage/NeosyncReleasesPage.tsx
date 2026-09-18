import { useState } from 'react';
import PageLayout from '../PageLayout';
import { useTranslation } from 'react-i18next';
import { basicMetaGetReleases } from '@/rest/metaAPI';
import { Loading } from '@/components/common/Loader/Loading';
import { useHandleServer } from '@/hooks/Server/useHandleServer';

const NeosyncReleasesPage = () => {
	const { t } = useTranslation();

	const currentIndex = 0;
	const [openIndex, setOpenIndex] = useState<number | null>(null);

	const lang = localStorage.getItem('lang') || 'ru';

	const { data: respMetaGetReleases, loading: respMetaGetReleasesLoading } = useHandleServer(['respMetaGetReleasesLoading'], basicMetaGetReleases);

	return (
		<PageLayout>
			<p className="font-medium text-[17px]">{t('message.version-history')}</p>

			<div className="mt-5 overflow-y-auto">
				{respMetaGetReleasesLoading ? (
					<Loading />
				) : (
					respMetaGetReleases?.releases.map((release, index) => {
						const isCurrent = index === currentIndex;
						const isOpen = isCurrent || openIndex === index;

						return (
							<div className="hover:bg-[#f6f6f6]" key={index}>
								<div
									className="flex items-center justify-between border-t border-[#e4e4e4] p-4 cursor-pointer"
									onClick={() => {
										if (!isCurrent) {
											setOpenIndex((prev) => (prev === index ? null : index));
										}
									}}
								>
									<div className="flex items-center gap-2 -mb-1">
										<p className="text-[16px]">
											{release.date} v{release.version}
										</p>
									</div>

									{isCurrent ? (
										<div className="border border-[#afb5c0] px-3 py-2 rounded">
											<p className="text-[#afb5c0] text-sm">{t('message.current-version')}</p>
										</div>
									) : (
										<div className="border border-[#395d95] px-3 py-2 rounded">
											<p className="text-[#395d95] text-sm">{isOpen ? t('label.hide') : t('label.more')}</p>
										</div>
									)}
								</div>

								{isOpen && (
									<div className="px-5 pb-5">
										<ul className="list-disc pl-8 -space-y-1 marker:text-[#49525f] marker:text-[20px]">
											{release.changes[lang].added &&
												release.changes[lang].added.map((ad, index) => (
													<li className="text-[14px]" key={`added_${index}`}>
														{ad}
													</li>
												))}
											{release.changes[lang].fixed &&
												release.changes[lang].fixed.map((fix, index) => (
													<li className="text-[14px]" key={`fixed_${index}`}>
														{fix}
													</li>
												))}
										</ul>
									</div>
								)}
							</div>
						);
					})
				)}
			</div>
		</PageLayout>
	);
};

export default NeosyncReleasesPage;
