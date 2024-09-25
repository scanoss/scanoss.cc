import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';

type UseQueryStateReturnType<T> = [T, Dispatch<SetStateAction<T>>];

const useQueryState = <T>(
  param: string,
  initialValue: T
): UseQueryStateReturnType<T> => {
  const [searchParams, setSearchParams] = useSearchParams();

  const paramValue = searchParams.get(param);
  let defaultValue: T;
  if (paramValue !== null) {
    defaultValue = paramValue as T;
  } else {
    defaultValue = initialValue;
  }

  const [value, setValue] = useState<T>(defaultValue);

  useEffect(() => {
    setSearchParams((prev) => {
      if (value) {
        prev.set(param, value.toString());
      } else prev.delete(param);
      return prev;
    });
  }, [value]);

  return [value, setValue];
};

export default useQueryState;
