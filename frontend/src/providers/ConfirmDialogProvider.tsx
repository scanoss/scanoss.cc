import { createContext, Dispatch, ReactNode, SetStateAction, useState } from 'react';

export type ConfirmContext = {
  isAsking: boolean;
  setIsAsking: Dispatch<SetStateAction<boolean>>;
  message?: ReactNode;
  setMessage: Dispatch<SetStateAction<ConfirmContext['message']>>;
  resolve?: ((value: boolean) => void) | undefined;
  setResolve: Dispatch<SetStateAction<ConfirmContext['resolve']>>;
};

export const confirmContext = createContext<ConfirmContext | null>(null);

export const ConfirmDialogProvider = ({ children }: { children: ReactNode }) => {
  const [isAsking, setIsAsking] = useState<ConfirmContext['isAsking']>(false);
  const [message, setMessage] = useState<ConfirmContext['message']>();
  const [resolve, setResolve] = useState<ConfirmContext['resolve']>();

  return (
    <confirmContext.Provider
      value={{
        isAsking,
        setIsAsking,
        message,
        setMessage,
        resolve,
        setResolve,
      }}
    >
      {children}
    </confirmContext.Provider>
  );
};
