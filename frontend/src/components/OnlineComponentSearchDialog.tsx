// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2025 SCANOSS.COM
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

import { useQuery } from '@tanstack/react-query';
import { Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';

import useDebounce from '@/hooks/useDebounce';

import { entities } from '../../wailsjs/go/models';
import { SearchComponents } from '../../wailsjs/go/service/ComponentServiceImpl';
import { ComponentSearchTable } from './ComponentSearchTable';
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandList } from './ui/command';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from './ui/dialog';
import { useToast } from './ui/use-toast';

const DEBOUNCE_QUERY_MS = 500;
const STATUS_CODE_OK = 1;

interface OnlineComponentSearchDialogProps {
  onOpenChange: (open: boolean) => void;
  searchTerm: string;
  onComponentSelect: (component: entities.SearchedComponent) => void;
}

export default function OnlineComponentSearchDialog({ onOpenChange, searchTerm, onComponentSelect }: OnlineComponentSearchDialogProps) {
  const { toast } = useToast();
  const [currentSearch, setCurrentSearch] = useState<string>(searchTerm);
  const debouncedSearchTerm = useDebounce<string>(currentSearch, DEBOUNCE_QUERY_MS);

  const {
    data: searchResults,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['components', debouncedSearchTerm],
    queryFn: async () => {
      const request: entities.ComponentSearchRequest = {
        search: debouncedSearchTerm.trim(),
        vendor: '',
        component: '',
        package: '',
        limit: 20,
        offset: 0,
      };

      if (!request.search) {
        return [];
      }

      const response = await SearchComponents(request);

      if (response.status.code !== STATUS_CODE_OK) {
        throw new Error(response.status.message);
      }

      return response.components;
    },
  });

  const handleComponentSelect = (component: entities.SearchedComponent) => {
    onComponentSelect(component);
    onOpenChange(false);
  };

  const errorMessage = error instanceof Error ? error.message : 'Failed to search for components online';

  useEffect(() => {
    if (error) {
      console.error('Error searching components:', error);
      toast({
        title: 'Search Error',
        description: errorMessage + '. Please try again.',
        variant: 'destructive',
      });
    }
  }, [error]);

  return (
    <Dialog open onOpenChange={onOpenChange}>
      <DialogContent className="p-4">
        <DialogHeader>
          <DialogTitle>Search Components Online</DialogTitle>
          <DialogDescription>Search for components in the SCANOSS database and select one to replace with</DialogDescription>
        </DialogHeader>

        <Command shouldFilter={false}>
          <CommandInput placeholder="Search for components..." defaultValue={searchTerm} onValueChange={(newValue) => setCurrentSearch(newValue)} />
          <CommandList>
            {isLoading && (
              <div className="flex items-center justify-center py-8">
                <Loader2 className="h-6 w-6 animate-spin" />
                <span className="ml-2 text-sm text-muted-foreground">Searching for &ldquo;{currentSearch}&rdquo;...</span>
              </div>
            )}

            {!isLoading && error && (
              <CommandEmpty>
                <div className="text-center text-destructive">
                  <p className="font-medium">Search Failed</p>
                  <p className="text-sm text-muted-foreground">{errorMessage}</p>
                  <p className="mt-1 text-xs text-muted-foreground">Please check your connection and try again.</p>
                </div>
              </CommandEmpty>
            )}

            {!isLoading && !error && debouncedSearchTerm && searchResults?.length === 0 && (
              <CommandEmpty>No components found for &ldquo;{debouncedSearchTerm}&rdquo;. Try a different search term.</CommandEmpty>
            )}

            {!isLoading && !error && !debouncedSearchTerm && <CommandEmpty>Enter a search term to find components online.</CommandEmpty>}

            {!isLoading && searchResults && searchResults.length > 0 && (
              <CommandGroup heading={`Found ${searchResults.length} components`}>
                <ComponentSearchTable components={searchResults} onComponentSelect={handleComponentSelect} />
              </CommandGroup>
            )}
          </CommandList>
        </Command>
      </DialogContent>
    </Dialog>
  );
}
