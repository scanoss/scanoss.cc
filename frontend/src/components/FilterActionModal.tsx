// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2026 SCANOSS.COM
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
import { Check, ChevronsUpDown, File, Folder, Package, Plus, Search } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';

import { useDialogRegistration } from '@/contexts/DialogStateContext';
import useEnvironment from '@/hooks/useEnvironment';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { getShortcutDisplay, KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { cn, getPathSegments } from '@/lib/utils';
import { FilterAction, filterActionLabelMap } from '@/modules/components/domain';
import { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import { entities } from '../../wailsjs/go/models';
import { GetDeclaredComponents } from '../../wailsjs/go/service/ComponentServiceImpl';
import { GetLicensesByPurl } from '../../wailsjs/go/service/LicenseServiceImpl';
import NewComponentDialog from './NewComponentDialog';
import PathBreadcrumb from './PathBreadcrumb';
import OnlineComponentSearchDialog from './OnlineComponentSearchDialog';
import SelectLicenseList from './SelectLicenseList';
import ShortcutBadge from './ShortcutBadge';
import { Button } from './ui/button';
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator } from './ui/command';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Label } from './ui/label';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { ScrollArea } from './ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Textarea } from './ui/textarea';
import { useToast } from './ui/use-toast';

function buildPath(segments: string[], index: number): string {
  const path = segments.slice(0, index + 1).join('/');
  const isFile = index === segments.length - 1;
  return isFile ? path : path + '/';
}

interface FilterActionModalProps {
  action: FilterAction;
  filePath: string;
  purl: string;
  affectedFilesCount: number;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: (args: OnFilterComponentArgs) => void;
}

export default function FilterActionModal({
  action,
  filePath,
  purl,
  affectedFilesCount,
  open,
  onOpenChange,
  onConfirm,
}: FilterActionModalProps) {
  const { toast } = useToast();
  const { modifierKey } = useEnvironment();
  const submitButtonRef = useRef<HTMLButtonElement>(null);

  // Tab state: 'path' or 'component'
  const [activeTab, setActiveTab] = useState<'path' | 'component'>('path');

  // Path selection state
  const segments = useMemo(() => getPathSegments(filePath), [filePath]);
  const [selectedPathIndex, setSelectedPathIndex] = useState(segments.length - 1);

  // Comment state
  const [comment, setComment] = useState('');

  // License state (for Include action)
  const [license, setLicense] = useState<string>();
  const [licenseKey, setLicenseKey] = useState(0);

  // Replace component state
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [newComponentDialogOpen, setNewComponentDialogOpen] = useState(false);
  const [onlineSearchDialogOpen, setOnlineSearchDialogOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');
  const [declaredComponents, setDeclaredComponents] = useState<entities.DeclaredComponent[]>([]);
  const [selectedComponent, setSelectedComponent] = useState<entities.DeclaredComponent | null>(null);
  const [matchedLicenses, setMatchedLicenses] = useState<entities.License[]>([]);

  const isIncludeAction = action === FilterAction.Include;
  const isReplaceAction = action === FilterAction.Replace;

  // Determine target type for dynamic text
  const getTargetType = () => {
    if (activeTab === 'component') return 'component';
    return selectedPathIndex < segments.length - 1 ? 'folder' : 'file';
  };
  const targetType = getTargetType();

  useDialogRegistration('filter-action', open);

  // Fetch declared components for Replace action
  const { data: componentsData, error: componentsError } = useQuery({
    queryKey: ['declaredComponents'],
    queryFn: GetDeclaredComponents,
    enabled: isReplaceAction && open,
  });

  useEffect(() => {
    if (componentsData) {
      setDeclaredComponents(componentsData);
    }
  }, [componentsData]);

  useEffect(() => {
    if (componentsError) {
      toast({
        title: 'Error',
        description: 'Error fetching declared components',
        variant: 'destructive',
      });
    }
  }, [componentsError, toast]);

  // Reset state when modal opens
  useEffect(() => {
    if (open) {
      setActiveTab('path');
      setSelectedPathIndex(segments.length - 1);
      setComment('');
      setLicense(undefined);
      setLicenseKey((k) => k + 1);
      setSelectedComponent(null);
      setSearchValue('');
    }
  }, [open, segments.length]);

  const getMatchedLicenses = async (componentPurl: string) => {
    setMatchedLicenses([]);
    try {
      const { component } = await GetLicensesByPurl({ purl: componentPurl });
      if (!component || !component.licenses) {
        setMatchedLicenses([]);
        return;
      }
      const licenses: entities.License[] = component.licenses.map((lic) => ({
        name: lic.full_name,
        licenseId: lic.id,
        reference: `https://spdx.org/licenses/${lic.id}.html`,
      }));
      setMatchedLicenses(licenses);
    } catch {
      setMatchedLicenses([]);
    }
  };

  const handleComponentSelected = (component: entities.DeclaredComponent) => {
    const alreadyExists = declaredComponents.some((c) => c.purl === component.purl);
    if (!alreadyExists) {
      setDeclaredComponents((prev) => [...prev, component]);
    }
    setSelectedComponent(component);
    setLicense(undefined);
    setLicenseKey((k) => k + 1);
    setNewComponentDialogOpen(false);
  };

  const handleSearchOnline = () => {
    setOnlineSearchDialogOpen(true);
    setPopoverOpen(false);
  };

  const handleConfirm = () => {
    const selectedPath = buildPath(segments, selectedPathIndex);
    const isFileSelected = selectedPathIndex === segments.length - 1;

    // Determine filterBy based on tab and path selection
    let filterBy: 'by_file' | 'by_purl' | 'by_folder';
    if (activeTab === 'component') {
      filterBy = 'by_purl';
    } else if (isFileSelected) {
      filterBy = 'by_file';
    } else {
      filterBy = 'by_folder';
    }

    // Validate for Replace action
    if (isReplaceAction && !selectedComponent) {
      toast({
        title: 'Error',
        description: 'Please select a component to replace with.',
        variant: 'destructive',
      });
      return;
    }

    onConfirm({
      action,
      filterBy,
      comment: comment || undefined,
      license: license || undefined,
      replaceWith: selectedComponent?.purl,
      folderPath: filterBy === 'by_folder' ? selectedPath : undefined,
    });

    onOpenChange(false);
  };

  const ref = useKeyboardShortcut(KEYBOARD_SHORTCUTS.confirm.keys, () => submitButtonRef.current?.click(), {
    enableOnFormTags: true,
    enabled: open,
  });

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.confirm.keys, handleSearchOnline, {
    enableOnFormTags: true,
    enabled: popoverOpen,
  });

  const actionLabel = filterActionLabelMap[action];

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent ref={ref} tabIndex={-1} className="p-4 max-w-lg">
          <DialogHeader>
            <DialogTitle>{actionLabel}</DialogTitle>
            <DialogDescription>
              {isReplaceAction
                ? `Replace the selected ${targetType} with another component.`
                : `${actionLabel} the selected ${targetType}.`}
            </DialogDescription>
          </DialogHeader>

          <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as 'path' | 'component')}>
            <TabsList className="self-start">
              <TabsTrigger value="path">Path</TabsTrigger>
              <TabsTrigger value="component">Component</TabsTrigger>
            </TabsList>
            <TabsContent value="path" className="min-h-[80px]">
              <div className="flex flex-col gap-2">
                <div className="rounded-md border bg-muted/30 px-3 py-2">
                  <PathBreadcrumb filePath={filePath} selectedIndex={selectedPathIndex} onSelect={setSelectedPathIndex} />
                </div>
                <span className="inline-flex items-center gap-1 self-start rounded-md border bg-muted/50 px-2 py-1 text-xs text-muted-foreground">
                  {selectedPathIndex < segments.length - 1 ? (
                    <>
                      <Folder className="h-3 w-3" />
                      folder
                    </>
                  ) : (
                    <>
                      <File className="h-3 w-3" />
                      file
                    </>
                  )}
                </span>
              </div>
            </TabsContent>
            <TabsContent value="component" className="min-h-[80px]">
              <div className="flex flex-col gap-2">
                <div className="flex items-center gap-2 rounded-md border bg-muted/30 px-3 py-2 font-mono text-sm">
                  <Package className="h-4 w-4 text-muted-foreground" />
                  <span>{purl}</span>
                </div>
                <span className="inline-flex items-center gap-1 self-start rounded-md border bg-muted/50 px-2 py-1 text-xs text-muted-foreground">
                  <Package className="h-3 w-3" />
                  {affectedFilesCount} pending {affectedFilesCount === 1 ? 'file' : 'files'}
                </span>
              </div>
            </TabsContent>
          </Tabs>

          {/* License selector - Include action only */}
          {isIncludeAction && (
            <div className="flex flex-col gap-2">
              <Label>
                License <span className="text-xs font-normal text-muted-foreground">(optional)</span>
              </Label>
              <SelectLicenseList key={licenseKey} onSelect={setLicense} />
            </div>
          )}

          {/* Component selector - Replace action only */}
          {isReplaceAction && (
            <div className="flex flex-col gap-2">
              <Label>Replace with</Label>
              <Popover open={popoverOpen} onOpenChange={setPopoverOpen} modal>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    role="combobox"
                    className={cn('justify-between', !selectedComponent && 'text-muted-foreground')}
                  >
                    {selectedComponent ? selectedComponent.name : 'Select component'}
                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="p-0">
                  <Command>
                    <CommandInput placeholder="Search component..." value={searchValue} onValueChange={setSearchValue} />
                    <CommandList>
                      <CommandEmpty>No components found.</CommandEmpty>
                      <CommandGroup>
                        <CommandItem
                          onSelect={() => {
                            setNewComponentDialogOpen(true);
                            setPopoverOpen(false);
                          }}
                        >
                          <Plus className="mr-2 h-4 w-4" />
                          Add new component
                        </CommandItem>
                        {searchValue && (
                          <CommandItem onSelect={handleSearchOnline} className="flex items-center justify-between gap-2">
                            <div className="flex items-center gap-2">
                              <Search className="h-3 w-3" />
                              Search &ldquo;{searchValue}&rdquo; online
                            </div>
                            <ShortcutBadge shortcut={getShortcutDisplay(KEYBOARD_SHORTCUTS.confirm.keys, modifierKey.label)[0]} />
                          </CommandItem>
                        )}
                      </CommandGroup>
                      <CommandSeparator />
                      <CommandGroup>
                        <ScrollArea className="h-52">
                          {declaredComponents?.map((component) => (
                            <CommandItem
                              key={component.purl}
                              onSelect={() => {
                                setSelectedComponent(component);
                                setLicense(undefined);
                                setLicenseKey((k) => k + 1);
                                setPopoverOpen(false);
                              }}
                            >
                              <Check
                                className={cn(
                                  'mr-2 h-4 w-4',
                                  component.purl === selectedComponent?.purl ? 'opacity-100' : 'opacity-0'
                                )}
                              />
                              <div>
                                <p>{component.name}</p>
                                <p className="text-xs text-muted-foreground">{component.purl}</p>
                              </div>
                            </CommandItem>
                          ))}
                        </ScrollArea>
                      </CommandGroup>
                    </CommandList>
                  </Command>
                </PopoverContent>
              </Popover>
              {selectedComponent && (
                <p className="text-xs text-muted-foreground font-mono">{selectedComponent.purl}</p>
              )}

              {/* License for replacement component */}
              <div className="flex flex-col gap-2 mt-2">
                <Label>
                  License <span className="text-xs font-normal text-muted-foreground">(optional)</span>
                </Label>
                <SelectLicenseList key={licenseKey} onSelect={setLicense} matchedLicenses={matchedLicenses} />
              </div>
            </div>
          )}

          {/* Comment textarea - always shown */}
          <div className="flex flex-col gap-2">
            <Label>
              Comment <span className="text-xs font-normal text-muted-foreground">(optional)</span>
            </Label>
            <Textarea value={comment} onChange={(e) => setComment(e.target.value)} placeholder="Add a comment..." />
          </div>

          <DialogFooter>
            <Button variant="ghost" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button ref={submitButtonRef} onClick={handleConfirm}>
              {actionLabel} <ShortcutBadge shortcut={getShortcutDisplay(KEYBOARD_SHORTCUTS.confirm.keys, modifierKey.label)[0]} />
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {newComponentDialogOpen && (
        <NewComponentDialog
          onOpenChange={() => setNewComponentDialogOpen(false)}
          onCreated={async (c) => {
            handleComponentSelected(c);
            await getMatchedLicenses(c.purl);
          }}
        />
      )}
      {onlineSearchDialogOpen && (
        <OnlineComponentSearchDialog
          onOpenChange={() => setOnlineSearchDialogOpen(false)}
          searchTerm={searchValue}
          onComponentSelect={async (c) => {
            handleComponentSelected(c);
            await getMatchedLicenses(c.purl);
          }}
        />
      )}
    </>
  );
}
