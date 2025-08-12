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
import { Loader2, Search } from 'lucide-react';
import { useEffect, useState } from 'react';

import useDebounce from '@/hooks/useDebounce';

import { entities } from '../../wailsjs/go/models';
import { SearchComponents } from '../../wailsjs/go/service/ComponentServiceImpl';
import { ComponentSearchTable } from './ComponentSearchTable';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from './ui/dialog';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { ScrollArea } from './ui/scroll-area';
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from './ui/select';
import { useToast } from './ui/use-toast';

const DEBOUNCE_QUERY_MS = 500;
const STATUS_CODE_OK = 1;

const catalogPackages = [
  'angular',
  'apache',
  'apple',
  'bitbucket',
  'cargo',
  'composer',
  'cpan',
  'deb',
  'drupal',
  'eclipse',
  'gem',
  'gitee',
  'github',
  'gitlab',
  'gnome',
  'gnu',
  'golang',
  'googlesource',
  'isc',
  'java2s',
  'jquery',
  'kernel',
  'maven',
  'mozilla',
  'nasm',
  'nmap',
  'npm',
  'nuget',
  'postgresql',
  'pypi',
  'rpm',
  'slf4j',
  'sourceforge',
  'stackoverflow',
  'sudo',
  'videolan',
  'wordpress',
  'zlib',
];

const DEFAULT_PACKAGE = 'github';

interface OnlineComponentSearchDialogProps {
  onOpenChange: (open: boolean) => void;
  searchTerm: string;
  onComponentSelect: (component: entities.DeclaredComponent) => void;
}

export default function OnlineComponentSearchDialog({ onOpenChange, searchTerm, onComponentSelect }: OnlineComponentSearchDialogProps) {
  const { toast } = useToast();
  const [currentSearch, setCurrentSearch] = useState<string>(searchTerm);
  const [selectedPackage, setSelectedPackage] = useState<string>(DEFAULT_PACKAGE);
  const debouncedSearchTerm = useDebounce<string>(currentSearch, DEBOUNCE_QUERY_MS);

  const {
    data: searchResults,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['components', debouncedSearchTerm, selectedPackage],
    queryFn: async () => {
      const request: entities.ComponentSearchRequest = {
        search: debouncedSearchTerm.trim(),
        package: selectedPackage,
      };

      if (!request.search) {
        return [];
      }

      const response = await SearchComponents(request);

      return response.components;
    },
  });

  const handleComponentSelect = (component: entities.SearchedComponent) => {
    const newComponent: entities.DeclaredComponent = {
      name: component.component,
      purl: component.purl,
    };
    onComponentSelect(newComponent);
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
      <DialogContent className="max-h-[calc(80dvh)] max-w-[calc(60dvw)] p-4 sm:max-w-[calc(90vw)] md:max-w-[calc(60dvw)]">
        <DialogHeader>
          <DialogTitle>Online Component Search</DialogTitle>
        </DialogHeader>

        <div className="flex flex-col gap-4">
          <div className="flex items-center gap-4">
            <div className="flex-1 space-y-2">
              <Label htmlFor="search">Search</Label>
              <div className="relative">
                <Input
                  className="pl-8"
                  id="search"
                  onChange={(e) => setCurrentSearch(e.target.value)}
                  placeholder="Search for components..."
                  value={currentSearch}
                />
                <Search className="absolute left-3 top-1/2 h-3 w-3 -translate-y-1/2 transform text-muted-foreground" />
              </div>
            </div>
            <div className="flex items-center gap-2">
              <div className="space-y-2">
                <Label htmlFor="package">Package</Label>
                <Select value={selectedPackage} onValueChange={(value) => setSelectedPackage(value)}>
                  <SelectTrigger className="w-[220px] capitalize">
                    <SelectValue placeholder="Select a package" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      {catalogPackages.map((catalogPackage) => (
                        <SelectItem key={catalogPackage} value={catalogPackage} className="capitalize">
                          {catalogPackage}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
            </div>
          </div>
          {isLoading && (
            <div className="flex items-center justify-center">
              <Loader2 className="h-6 w-6 animate-spin" />
              <span className="ml-2 text-sm text-muted-foreground">Searching for &ldquo;{currentSearch}&rdquo;...</span>
            </div>
          )}

          {!isLoading && error && (
            <div className="flex items-center justify-center">
              <div className="text-center text-destructive">
                <p className="font-medium">Search Failed</p>
                <p className="text-sm text-muted-foreground">{errorMessage}</p>
                <p className="mt-1 text-xs text-muted-foreground">Please check your connection and try again.</p>
              </div>
            </div>
          )}

          {!isLoading && !error && debouncedSearchTerm && searchResults?.length === 0 && (
            <div className="flex items-center justify-center text-muted-foreground">
              <p>No components found for &ldquo;{debouncedSearchTerm}&rdquo;. Try a different search term.</p>
            </div>
          )}

          {!isLoading && !error && !debouncedSearchTerm && (
            <div className="flex items-center justify-center text-muted-foreground">Enter a search term to find components online.</div>
          )}

          {!isLoading && searchResults && searchResults.length > 0 && (
            <div className="flex flex-col gap-4 overflow-auto">
              <p className="text-sm text-muted-foreground">Found {searchResults.length} components</p>
              <ScrollArea className="max-h-[300px]">
                <ComponentSearchTable components={searchResults} onComponentSelect={handleComponentSelect} />
              </ScrollArea>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
