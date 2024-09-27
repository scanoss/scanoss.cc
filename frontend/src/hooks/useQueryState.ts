import { useCallback, useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';

type Serializable = string | number | boolean | null | undefined;

type SerializeFunction<T> = (value: T) => string;
type DeserializeFunction<T> = (value: string) => T;

const defaultSerialize = <T extends Serializable>(value: T): string => {
  return JSON.stringify(value);
};

const defaultDeserialize = <T extends Serializable>(value: string): T => {
  try {
    return JSON.parse(value) as T;
  } catch {
    return value as T;
  }
};

const useQueryState = <T extends Serializable>(
  key: string,
  initialValue: T,
  options?: {
    serialize?: SerializeFunction<T>;
    deserialize?: DeserializeFunction<T>;
  }
): [T, (newValue: T) => void] => {
  const [searchParams, setSearchParams] = useSearchParams();
  const serialize = options?.serialize ?? defaultSerialize;
  const deserialize = options?.deserialize ?? defaultDeserialize;

  const [value, setValue] = useState<T>(() => {
    const paramValue = searchParams.get(key);
    if (paramValue === null) {
      return initialValue;
    }
    try {
      return deserialize(paramValue);
    } catch {
      return initialValue;
    }
  });

  useEffect(() => {
    const paramValue = searchParams.get(key);
    if (paramValue !== null) {
      try {
        const deserializedValue = deserialize(paramValue);
        if (deserializedValue !== value) {
          setValue(deserializedValue);
        }
      } catch {
        if (value !== initialValue) {
          setValue(initialValue);
        }
      }
    } else if (value !== initialValue) {
      setValue(initialValue);
    }
  }, [searchParams, key, initialValue, value, deserialize]);

  const updateValue = useCallback(
    (newValue: T) => {
      setValue(newValue);
      setSearchParams((params) => {
        if (!newValue) {
          params.delete(key);
        } else {
          params.set(key, serialize(newValue));
        }
        return params;
      });
    },
    [key, setSearchParams, serialize]
  );

  return [value, updateValue];
};

export default useQueryState;
