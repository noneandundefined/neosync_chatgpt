import i18next from 'i18next';
import { toast } from 'react-toastify';

export const copyToBuffer = async (data: string) => {
	if (!data) return;

	try {
		if (navigator.clipboard && window.isSecureContext) {
			await navigator.clipboard.writeText(data);
		} else {
			const textArea = document.createElement('textarea');
			textArea.value = data;
			textArea.style.position = 'fixed';
			textArea.style.left = '-9999px';
			document.body.appendChild(textArea);
			textArea.focus();
			textArea.select();
			document.execCommand('copy');
			document.body.removeChild(textArea);
		}

		toast.success(i18next.t('message.copy-success'));
	} catch {
		toast.error(i18next.t('message.copy-error'));
	}
};
