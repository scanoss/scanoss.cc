import { ReactNode, useContext, useEffect } from 'react';

import {
  ConfirmContext,
  confirmContext,
} from '@/providers/ConfirmDialogProvider';

export function useConfirm() {
  // @ts-expect-error - Check if context is not null
  const context = useContext<ConfirmContext>(confirmContext);

  if (!context) {
    throw new Error(
      "'useConfirm' is being used outside of ConfirmDialogProvider"
    );
  }

  const {
    message,
    setMessage,
    resolve,
    setResolve,
    setIsAsking,
    isAsking,
    setOnPersistDecision,
    onPersistDecision,
  } = context;

  const ask = async (
    msg: ReactNode,
    onPersistDecision?: () => void
  ): Promise<boolean> => {
    return new Promise((resolve) => {
      setMessage(msg);
      setIsAsking(true);
      setResolve(() => (value: boolean) => {
        resolve(value);
        setOnPersistDecision(onPersistDecision);
      });
    });
  };

  const confirm = () => {
    resolve?.(true);
    setIsAsking(false);
  };

  const deny = () => {
    resolve?.(false);
    setIsAsking(false);
  };

  useEffect(() => {
    if (!isAsking) {
      setTimeout(() => setMessage(undefined), 300);
    }
  }, [isAsking]);

  return { message, isAsking, ask, confirm, deny, onPersistDecision };
}
