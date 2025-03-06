import { useQueryClient } from '@tanstack/react-query';
import { Loader2, Save, XCircle } from 'lucide-react';
import { useEffect, useState } from 'react';

import useSkipPatternStore from '@/hooks/useSkipPatternStore';
import useConfigStore from '@/stores/useConfigStore';

import { CommitStagedSkipPatterns, DiscardStagedSkipPatterns, HasStagedChanges } from '../../wailsjs/go/service/ScanossSettingsServiceImp';
import { Button } from './ui/button';
import { useToast } from './ui/use-toast';

export default function SaveButton() {
  const queryClient = useQueryClient();
  const { toast } = useToast();
  const scanRoot = useConfigStore((state) => state.scanRoot);
  const [isLoading, setIsLoading] = useState(false);
  const { hasUnsavedChanges, setHasUnsavedChanges, clearNodesWithUnsavedChanges } = useSkipPatternStore();

  useEffect(() => {
    const checkUnsavedChanges = async () => {
      try {
        const hasChanges = await HasStagedChanges();
        setHasUnsavedChanges(hasChanges);
      } catch (error) {
        console.error('Error checking for unsaved changes', error);
      }
    };

    checkUnsavedChanges();

    const interval = setInterval(checkUnsavedChanges, 5000);
    return () => clearInterval(interval);
  }, [setHasUnsavedChanges]);

  const commitStagedSkipPatterns = async () => {
    try {
      setIsLoading(true);
      await CommitStagedSkipPatterns();

      toast({
        title: 'Success',
        description: 'Scanning skip patterns saved successfully',
      });

      queryClient.invalidateQueries({ queryKey: ['resultsTree', scanRoot] });
      setHasUnsavedChanges(false);
      clearNodesWithUnsavedChanges();
    } catch (error) {
      console.error('Error saving skip patterns', error);
      toast({
        title: 'Error',
        description: 'Failed to save scanning skip patterns',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  const discardStagedSkipPatterns = async () => {
    try {
      setIsLoading(true);
      await DiscardStagedSkipPatterns();

      toast({
        title: 'Changes Discarded',
        description: 'Scanning skip pattern changes have been discarded',
      });

      queryClient.invalidateQueries({ queryKey: ['resultsTree', scanRoot] });
      setHasUnsavedChanges(false);
      clearNodesWithUnsavedChanges();
    } catch (error) {
      console.error('Error discarding skip patterns', error);
      toast({
        title: 'Error',
        description: 'Failed to discard scanning skip patterns',
        variant: 'destructive',
      });
    } finally {
      setIsLoading(false);
    }
  };

  if (!hasUnsavedChanges) {
    return null;
  }

  return (
    <div className="fixed bottom-4 right-4 z-50 flex gap-2 rounded-lg border border-border bg-background p-4 shadow-lg animate-in fade-in slide-in-from-bottom-5">
      <div className="flex flex-col gap-2">
        <div className="mb-1 text-sm font-medium">You have unsaved changes</div>
        <div className="flex gap-2">
          <Button onClick={commitStagedSkipPatterns} disabled={isLoading} variant="default" className="w-full">
            {isLoading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Save className="mr-2 h-4 w-4" />}
            Save Changes
          </Button>
          <Button onClick={discardStagedSkipPatterns} disabled={isLoading} variant="secondary" className="w-full">
            <XCircle className="mr-2 h-4 w-4" />
            Discard
          </Button>
        </div>
      </div>
    </div>
  );
}
