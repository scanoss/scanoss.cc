import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Loader2, Save, XCircle } from 'lucide-react';

import useConfigStore from '@/stores/useConfigStore';

import {
  CommitStagedScanningSkipPatterns,
  DiscardStagedScanningSkipPatterns,
  HasStagedScanningSkipPatternChanges,
} from '../../wailsjs/go/service/ScanossSettingsServiceImp';
import { Button } from './ui/button';
import { useToast } from './ui/use-toast';

export default function SaveSkipSettingsButton() {
  const { toast } = useToast();
  const queryClient = useQueryClient();

  const scanRoot = useConfigStore((state) => state.scanRoot);

  const { data: hasUnsavedChanges } = useQuery({
    queryKey: ['hasStagedScanningSkipPatternChanges', scanRoot],
    queryFn: HasStagedScanningSkipPatternChanges,
    refetchInterval: 3000,
    refetchIntervalInBackground: false,
    retry: 3,
  });

  const { mutate: commitStagedSkipPatterns, isPending: isCommitting } = useMutation({
    mutationFn: async () => {
      await CommitStagedScanningSkipPatterns();
    },
    onSuccess: () => {
      toast({
        title: 'Success',
        description: 'Scanning skip patterns saved successfully',
      });
      queryClient.invalidateQueries({ queryKey: ['resultsTree', scanRoot] });
      queryClient.invalidateQueries({ queryKey: ['hasStagedScanningSkipPatternChanges', scanRoot] });
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to save scanning skip patterns',
        variant: 'destructive',
      });
    },
  });

  const { mutate: discardStagedSkipPatterns, isPending: isDiscarding } = useMutation({
    mutationFn: async () => {
      await DiscardStagedScanningSkipPatterns();
    },
    onSuccess: () => {
      toast({
        title: 'Changes Discarded',
        description: 'Scanning skip pattern changes have been discarded',
      });
      queryClient.invalidateQueries({ queryKey: ['resultsTree', scanRoot] });
      queryClient.invalidateQueries({ queryKey: ['hasStagedScanningSkipPatternChanges', scanRoot] });
    },
    onError: () => {
      toast({
        title: 'Error',
        description: 'Failed to discard scanning skip patterns',
        variant: 'destructive',
      });
    },
  });

  if (!hasUnsavedChanges) {
    return null;
  }

  return (
    <div className="fixed bottom-4 right-4 z-50 flex gap-2 rounded-lg border border-border bg-background p-4 shadow-lg animate-in fade-in slide-in-from-bottom-5">
      <div className="flex flex-col gap-2">
        <div className="mb-1 text-sm font-medium">You have unsaved changes</div>
        <div className="flex gap-2">
          <Button onClick={() => commitStagedSkipPatterns()} disabled={isCommitting} variant="default" className="w-full">
            {isCommitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Save className="mr-2 h-4 w-4" />}
            Save Changes
          </Button>
          <Button onClick={() => discardStagedSkipPatterns()} disabled={isDiscarding} variant="secondary" className="w-full">
            {isDiscarding ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <XCircle className="mr-2 h-4 w-4" />}
            Discard
          </Button>
        </div>
      </div>
    </div>
  );
}
