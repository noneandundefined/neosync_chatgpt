import ModalMove from '@/components/Modal/ModalMove';
import { createContext, ReactNode, useCallback, useContext, useEffect, useState } from 'react';

interface ModalMoveItem {
	id: string;
	content: ReactNode;
}

interface ModalMoveContextType {
	open: (content: ReactNode) => void;
	close: (id: string) => void;
	closeLast: () => void;
	bringToFront: (id: string) => void;
}

const ModalMoveContext = createContext<ModalMoveContextType | null>(null);

export const useModalMoveContext = () => {
	const context = useContext(ModalMoveContext);
	if (!context) {
		throw new Error('useModalMoveContext must be used inside ModalMoveProvider');
	}

	return context;
};

const generateId = (): string => {
	return Date.now().toString(36) + Math.random().toString(36).substr(2, 9);
};

export const ModalMoveProvider = ({ children }: { children: ReactNode }) => {
	const [modals, setModals] = useState<ModalMoveItem[]>([]);

	const open = useCallback((content: ReactNode) => {
		const id = generateId();
		setModals((prev) => [...prev, { id, content }]);
	}, []);

	const close = useCallback((id: string) => {
		setModals((prev) => prev.filter((m) => m.id !== id));
	}, []);

	const closeLast = useCallback(() => {
		setModals((prev) => prev.slice(0, -1));
	}, []);

	const bringToFront = useCallback((id: string) => {
		setModals((prev) => {
			const index = prev.findIndex((m) => m.id === id);
			if (index === -1 || index === prev.length - 1) return prev;
			const modal = prev[index];
			const newModals = [...prev];
			newModals.splice(index, 1);
			newModals.push(modal);
			return newModals;
		});
	}, []);

	useEffect(() => {
		const handleKey = (e: KeyboardEvent) => {
			if (e.key === 'Escape') {
				closeLast();
			}
		};

		if (modals.length > 0) {
			document.addEventListener('keydown', handleKey);
		}

		return () => {
			document.removeEventListener('keydown', handleKey);
		};
	}, [modals.length, closeLast]);

	return (
		<ModalMoveContext.Provider value={{ open, close, closeLast, bringToFront }}>
			{children}

			{modals.map((modal, index) => (
				<ModalMove key={modal.id} open={true} width="700px" zIndex={1000 + index} onClick={() => bringToFront(modal.id)} onClose={() => close(modal.id)}>
					{modal.content}
				</ModalMove>
			))}
		</ModalMoveContext.Provider>
	);
};
