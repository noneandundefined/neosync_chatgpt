import Modal from '@/components/Modal/Modal';
import { useTranslation } from 'react-i18next';
import { ROUTES } from '@/constants/constants';
import { useModalContext } from '@/context/useModalContext';
import { useCallback, useContext, useEffect, useRef, useState } from 'react';
import { UNSAFE_NavigationContext, resolvePath, To, useLocation, useNavigate } from 'react-router-dom';
import { clearTemplateDraft, EVENT_TEMPLATE_DRAFT, hasTemplateDraft } from '@/utils/TemplateDraftUtils';
import ModalExistRemoveTemplate from '@/components/common/Configurations/ConfigurationTemplate/Modals/ModalExistRemoveTemplate';

const isTemplateEditorPath = (pathname: string) => pathname === ROUTES.CONFIGURATIONS_TEMPLATE_NEW || /^\/configurations\/templates\/\d+$/.test(pathname);

const isLeavingTemplatePage = (targetPathname: string, currentPathname: string) => {
	const resolved = resolvePath(targetPathname, currentPathname);
	return isTemplateEditorPath(currentPathname) && !isTemplateEditorPath(resolved.pathname);
};

const getTargetPathname = (to: To, currentPathname: string) => resolvePath(to, currentPathname).pathname;

export const useTemplateExitGuard = (templateId?: number) => {
	const { t } = useTranslation();

	const location = useLocation();
	const navigate = useNavigate();

	const { open, close } = useModalContext();

	const { navigator } = useContext(UNSAFE_NavigationContext);

	const modalOpenRef = useRef(false);
	const proceedRef = useRef<(() => void) | null>(null);
	const locationRef = useRef(location);
	const [hasDraft, setHasDraft] = useState(() => hasTemplateDraft(templateId));

	useEffect(() => {
		const syncDraftState = () => setHasDraft(hasTemplateDraft(templateId));

		syncDraftState();
		window.addEventListener(EVENT_TEMPLATE_DRAFT, syncDraftState);

		return () => window.removeEventListener(EVENT_TEMPLATE_DRAFT, syncDraftState);
	}, [templateId]);

	useEffect(() => {
		locationRef.current = location;
	}, [location]);

	const showExitModal = useCallback(
		(proceed: () => void) => {
			if (modalOpenRef.current) return;

			modalOpenRef.current = true;
			proceedRef.current = proceed;

			open(
				<Modal title={t('message.exit-configuration-template-title')}>
					<ModalExistRemoveTemplate
						onClose={() => {
							modalOpenRef.current = false;
							proceedRef.current = null;
							close();
						}}
						onConfirm={() => {
							clearTemplateDraft(templateId);
							modalOpenRef.current = false;
							const next = proceedRef.current;
							proceedRef.current = null;
							close();
							next?.();
						}}
					/>
				</Modal>
			);
		},
		[open, close, t, templateId]
	);

	useEffect(() => {
		const handleBeforeUnload = (e: BeforeUnloadEvent) => {
			if (!hasTemplateDraft(templateId)) return;

			e.preventDefault();
			e.returnValue = '';
		};

		window.addEventListener('beforeunload', handleBeforeUnload);
		return () => window.removeEventListener('beforeunload', handleBeforeUnload);
	}, [templateId]);

	useEffect(() => {
		const handleClick = (e: MouseEvent) => {
			if (!hasTemplateDraft(templateId)) return;

			const anchor = (e.target as HTMLElement).closest('a[href]');
			if (!anchor || anchor.getAttribute('target') === '_blank') return;

			const href = anchor.getAttribute('href');
			if (!href || href.startsWith('#') || href.startsWith('mailto:') || href.startsWith('tel:')) return;

			const url = new URL(href, window.location.origin);
			if (!isLeavingTemplatePage(url.pathname, locationRef.current.pathname)) return;

			e.preventDefault();
			e.stopPropagation();

			const target = url.pathname + url.search + url.hash;
			showExitModal(() => navigate(target));
		};

		document.addEventListener('click', handleClick, true);
		return () => document.removeEventListener('click', handleClick, true);
	}, [navigate, showExitModal, templateId]);

	useEffect(() => {
		const originalPush = navigator.push;
		const originalReplace = navigator.replace;

		navigator.push = ((to: To, state?: unknown) => {
			const currentPathname = locationRef.current.pathname;

			if (!hasTemplateDraft(templateId) || !isLeavingTemplatePage(getTargetPathname(to, currentPathname), currentPathname)) {
				return originalPush.call(navigator, to, state);
			}

			showExitModal(() => originalPush.call(navigator, to, state));
		}) as typeof navigator.push;

		navigator.replace = ((to: To, state?: unknown) => {
			const currentPathname = locationRef.current.pathname;

			if (!hasTemplateDraft(templateId) || !isLeavingTemplatePage(getTargetPathname(to, currentPathname), currentPathname)) {
				return originalReplace.call(navigator, to, state);
			}

			showExitModal(() => originalReplace.call(navigator, to, state));
		}) as typeof navigator.replace;

		return () => {
			navigator.push = originalPush;
			navigator.replace = originalReplace;
		};
	}, [navigator, showExitModal, templateId]);

	useEffect(() => {
		if (!hasDraft) return;

		const trapUrl = window.location.href;
		window.history.pushState({ templateExitGuard: true }, '', trapUrl);

		const handlePopState = () => {
			if (!hasTemplateDraft(templateId)) return;

			window.history.pushState({ templateExitGuard: true }, '', trapUrl);

			showExitModal(() => {
				window.history.go(-2);
			});
		};

		window.addEventListener('popstate', handlePopState);
		return () => window.removeEventListener('popstate', handlePopState);
	}, [hasDraft, location.pathname, location.search, showExitModal, templateId]);
};

export default useTemplateExitGuard;
