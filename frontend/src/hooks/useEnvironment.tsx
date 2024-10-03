import { useEffect, useState } from 'react';

import { Environment, EnvironmentInfo } from '../../wailsjs/runtime/runtime';

interface UseEnvironmentReturnType {
  environment: EnvironmentInfo | undefined;
}

export default function useEnvironment(): UseEnvironmentReturnType {
  const [environment, setEnvironment] = useState<EnvironmentInfo>();

  useEffect(() => {
    async function fetchEnvironment() {
      const environment = await Environment();

      setEnvironment(environment);
    }

    fetchEnvironment();
  }, []);

  return {
    environment,
  };
}
