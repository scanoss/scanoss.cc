import { FileCog } from 'lucide-react';

import SkipSettings from '@/components/SkipSettings';

interface SettingsOption {
  id: string;
  title: string;
  description: string;
  icon: React.ElementType;
  component: React.ReactNode;
}

export const SETTINGS_OPTIONS: SettingsOption[] = [
  {
    id: 'skip-settings',
    title: 'Skip Settings',
    description: 'Configure the files, directories and extensions to skip while scanning',
    icon: FileCog,
    component: <SkipSettings />,
  },
];
