import { Check, PackageMinus, Replace } from 'lucide-react';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { FilterAction } from '@/modules/components/domain';
import useComponentFilterStore from '@/modules/components/stores/useComponentFilterStore';

import FilterActionButton from './FilterActionButton';
import ReplaceComponentDialog from './ReplaceComponentDialog';

export default function FilterComponentActions() {
  const navigate = useNavigate();

  const [showReplaceComponentDialog, setShowReplaceComponentDialog] =
    useState(false);

  const onFilterComponent = useComponentFilterStore(
    (state) => state.onFilterComponent
  );

  const handleFilterComponent = async () => {
    const nextResultRoute = await onFilterComponent();

    if (nextResultRoute) {
      navigate(nextResultRoute);
    }
  };

  const handleReplaceComponent = () => {
    setShowReplaceComponentDialog(true);
  };

  return (
    <>
      <div className="flex gap-2 md:justify-center">
        <FilterActionButton
          action={FilterAction.Include}
          icon={<Check className="h-5 w-5 stroke-green-500" />}
          description="By including a file/component, you force the engine to consider it with priority in future scans."
          onAdd={handleFilterComponent}
        />
        <FilterActionButton
          action={FilterAction.Remove}
          description="Dismissing a file/component will exclude it from future scan results."
          icon={<PackageMinus className="h-5 w-5 stroke-red-500" />}
          onAdd={handleFilterComponent}
        />
        <FilterActionButton
          action={FilterAction.Replace}
          description="Dismissing a file/component will exclude it from future scan results."
          icon={<Replace className="h-5 w-5 stroke-yellow-500" />}
          onAdd={handleReplaceComponent}
        />
      </div>
      {showReplaceComponentDialog && <ReplaceComponentDialog />}
    </>
  );
}
