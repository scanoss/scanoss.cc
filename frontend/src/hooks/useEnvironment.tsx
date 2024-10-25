import { useEffect, useState } from 'react';

import { Environment, EnvironmentInfo } from '../../wailsjs/runtime/runtime';

interface UseEnvironmentReturnType {
  environment: EnvironmentInfo | undefined;
  isMac: boolean;
  modifierKey: {
    label: string;
    keyCode: string;
  };
}

export default function useEnvironment(): UseEnvironmentReturnType {
  const [environment, setEnvironment] = useState<EnvironmentInfo>();

  const isMac = environment?.platform === 'darwin';
  const modifierKey = isMac
    ? {
        label: '⌘',
        keyCode: 'Meta',
      }
    : {
        label: 'Ctrl',
        keyCode: 'Control',
      };

  useEffect(() => {
    async function fetchEnvironment() {
      const environment = await Environment();

      setEnvironment(environment);
    }

    fetchEnvironment();
  }, []);

  return {
    environment,
    isMac,
    modifierKey,
  };
}
