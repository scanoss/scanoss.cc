import { useContext, useEffect } from 'react';

import {
  inputPromptContext,
  InputPromptOptions,
} from '@/providers/InputPromptDialogProvider';

export function useInputPrompt() {
  const context = useContext(inputPromptContext);

  if (!context) {
    throw new Error(
      "'useInputPrompt' is being used outside of PromptDialogProvider"
    );
  }

  const {
    setOptions,
    setIsPrompting,
    setResolve,
    isPrompting,
    resolve,
    options,
  } = context;

  const prompt = (options: InputPromptOptions): Promise<string | undefined> => {
    return new Promise((resolve) => {
      setOptions(options);
      setIsPrompting(true);
      setResolve(() => (value: string | undefined) => {
        resolve(value);
        setIsPrompting(false);
      });
    });
  };

  const confirm = (value: string) => {
    resolve?.(value);
    setIsPrompting(false);
  };

  const cancel = () => {
    resolve?.(undefined);
    setIsPrompting(false);
  };

  useEffect(() => {
    if (!isPrompting) {
      setTimeout(() => setOptions(null), 300);
    }
  }, [isPrompting]);

  return { prompt, confirm, cancel, isPrompting, options };
}
