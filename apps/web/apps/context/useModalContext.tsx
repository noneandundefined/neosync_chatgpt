import { createContext, ReactNode, useCallback, useContext, useEffect, useState } from 'react';

interface ModalContextType {
	open: (node: ReactNode) => void;
	close: () => void;
}

const ModalContext = createContext<ModalContextType | null>(null);

export const useModalContext = () => {
	const context = useContext(ModalContext);
	if (!context) {
		throw new Error('useModalContext must be used inside ModalProvider');
	}
	return context;
};

export const ModalProvider = ({ children }: { children: ReactNode }) => {
	const [modalNode, setModalNode] = useState<ReactNode | null>(null);

	const open = useCallback((node: ReactNode) => {
		setModalNode(node);
	}, []);

	const close = useCallback(() => {
		setModalNode(null);
	}, []);

	useEffect(() => {
		document.body.style.overflow = modalNode ? 'hidden' : '';
	}, [modalNode]);

	// ESC
	useEffect(() => {
		const handleEsc = (e: KeyboardEvent) => {
			if (e.key === 'Escape') close();
		};

		if (modalNode) document.addEventListener('keydown', handleEsc);

		return () => document.removeEventListener('keydown', handleEsc);
	}, [modalNode, close]);

	return (
		<ModalContext.Provider value={{ open, close }}>
			{children}
			{modalNode}
		</ModalContext.Provider>
	);
};
