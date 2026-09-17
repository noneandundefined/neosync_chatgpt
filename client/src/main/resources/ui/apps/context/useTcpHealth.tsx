import { createContext, useContext } from 'react';

export const TcpHealthContext = createContext<boolean>(true);

export const useTcpHealth = () => useContext(TcpHealthContext);
