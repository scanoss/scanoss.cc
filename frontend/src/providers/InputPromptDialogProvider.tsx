import React, { createContext, ReactNode, useState } from 'react';

export interface InputPromptOptions {
  title: string;
  description?: string;
  confirmText?: string;
  cancelText?: string;
  input: {
    defaultValue?: string;
    type: 'textarea' | 'text';
  };
}

export type InputPromptContext = {
  isPrompting: boolean;
  setIsPrompting: React.Dispatch<React.SetStateAction<boolean>>;
  options: InputPromptOptions | null;
  setOptions: React.Dispatch<React.SetStateAction<InputPromptOptions | null>>;
  resolve: ((value: string | undefined) => void) | undefined;
  setResolve: React.Dispatch<
    React.SetStateAction<((value: string | undefined) => void) | undefined>
  >;
};

export const inputPromptContext = createContext<InputPromptContext | null>(
  null
);

export const InputPromptDialogProvider = ({
  children,
}: {
  children: ReactNode;
}) => {
  const [isPrompting, setIsPrompting] = useState(false);
  const [options, setOptions] = useState<InputPromptOptions | null>(null);
  const [resolve, setResolve] = useState<
    ((value: string | undefined) => void) | undefined
  >();

  return (
    <inputPromptContext.Provider
      value={{
        isPrompting,
        setIsPrompting,
        options,
        setOptions,
        resolve,
        setResolve,
      }}
    >
      {children}
    </inputPromptContext.Provider>
  );
};
