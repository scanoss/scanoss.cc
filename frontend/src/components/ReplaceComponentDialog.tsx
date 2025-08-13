// SPDX-License-Identifier: MIT
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
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

import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { Check, ChevronsUpDown, CircleAlert, Plus, Search } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { useDialogRegistration } from '@/contexts/DialogStateContext';
import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { cn } from '@/lib/utils';
import { FilterAction } from '@/modules/components/domain';
import { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import { entities } from '../../wailsjs/go/models';
import { GetDeclaredComponents } from '../../wailsjs/go/service/ComponentServiceImpl';
import { GetLicensesByPurl } from '../../wailsjs/go/service/LicenseServiceImpl';
import FilterByPurlList from './FilterByPurlList';
import NewComponentDialog from './NewComponentDialog';
import OnlineComponentSearchDialog from './OnlineComponentSearchDialog';
import SelectLicenseList from './SelectLicenseList';
import ShortcutBadge from './ShortcutBadge';
import { Alert, AlertDescription, AlertTitle } from './ui/alert';
import { Button } from './ui/button';
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList, CommandSeparator } from './ui/command';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from './ui/form';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { ScrollArea } from './ui/scroll-area';
import { Textarea } from './ui/textarea';
import { useToast } from './ui/use-toast';

const ReplaceComponentFormSchema = z.object({
  name: z.string().min(1, 'You must select a component.'),
  purl: z.string().min(1, 'You must select a component.'),
  comment: z.string().optional(),
  license: z.string().optional(),
});

interface ReplaceComponentDialogProps {
  onOpenChange: () => void;
  onReplaceComponent: (args: OnFilterComponentArgs) => void;
  withComment: boolean;
  filterBy: 'by_file' | 'by_purl';
}

export default function ReplaceComponentDialog({ onOpenChange, onReplaceComponent, withComment, filterBy }: ReplaceComponentDialogProps) {
  const { toast } = useToast();
  const submitButtonRef = useRef<HTMLButtonElement>(null);
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [newComponentDialogOpen, setNewComponentDialogOpen] = useState(false);
  const [onlineSearchDialogOpen, setOnlineSearchDialogOpen] = useState(false);
  const [searchValue, setSearchValue] = useState('');
  const [declaredComponents, setDeclaredComponents] = useState<entities.DeclaredComponent[]>([]);
  const [licenseKey, setLicenseKey] = useState(0);
  const [matchedLicenses, setMatchedLicenses] = useState<entities.License[]>([]);
  
  useDialogRegistration('replace-component', true);

  const form = useForm<z.infer<typeof ReplaceComponentFormSchema>>({
    resolver: zodResolver(ReplaceComponentFormSchema),
    defaultValues: {
      name: '',
      purl: '',
      comment: '',
      license: '',
    },
  });

  const selectedPurl = form.watch('purl');

  const { data, error } = useQuery({
    queryKey: ['declaredComponents'],
    queryFn: GetDeclaredComponents,
  });

  const resetLicense = () => {
    form.setValue('license', '');
    setLicenseKey((prevKey) => prevKey + 1); // This is a hack to reset the SelectLicenseList component
  };

  const onSubmit = (values: z.infer<typeof ReplaceComponentFormSchema>) => {
    const { comment, purl, license } = values;

    onReplaceComponent({
      filterBy,
      comment,
      license,
      action: FilterAction.Replace,
      replaceWith: purl,
    });
  };

  const onComponentSelected = (component: entities.DeclaredComponent) => {
    const alreadyExists = declaredComponents.some((c) => c.purl === component.purl);
    if (alreadyExists) {
      return toast({
        title: 'Warning',
        description: `A component with the same PURL already exists. Please select it from the list or enter a different PURL.`,
      });
    }
    setDeclaredComponents((prevState) => [...prevState, component]);
    form.setValue('purl', component.purl);
    form.setValue('name', component.name);
    resetLicense();
    setNewComponentDialogOpen(false);
  };

  const getMatchedLicenses = async (purl: string) => {
    try {
      const { component } = await GetLicensesByPurl({ purl });

      if (!component || !component.licenses) return;

      const matchedLicenses: entities.License[] = component.licenses.map((license) => ({
        name: license.full_name,
        licenseId: license.id,
        reference: `https://spdx.org/licenses/${license.id}.html`, // This is a workaround for the fact that the reference is not available in the response
      }));
      setMatchedLicenses(matchedLicenses);
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Error fetching matched licenses',
        variant: 'destructive',
      });
    }
  };

  const handleSearchOnline = () => {
    setOnlineSearchDialogOpen(true);
    setPopoverOpen(false);
  };

  useEffect(() => {
    if (data) {
      setDeclaredComponents(data);
    }
  }, [data]);

  useEffect(() => {
    if (error) {
      toast({
        title: 'Error',
        description: 'Error fetching declared components',
        variant: 'destructive',
      });
    }
  }, [error]);

  const ref = useKeyboardShortcut(KEYBOARD_SHORTCUTS.confirm.keys, () => submitButtonRef.current?.click(), {
    enableOnFormTags: true,
  });

  useKeyboardShortcut(KEYBOARD_SHORTCUTS.confirm.keys, () => handleSearchOnline(), {
    enableOnFormTags: true,
    enabled: popoverOpen,
  });

  return (
    <>
      <Form {...form}>
        <Dialog open onOpenChange={onOpenChange}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <DialogContent ref={ref} tabIndex={-1} className="p-4">
              <DialogHeader>
                <DialogTitle>Replace</DialogTitle>
                <DialogDescription>You can search for an existing component or manually enter a PURL</DialogDescription>
              </DialogHeader>

              <FormField
                control={form.control}
                name="purl"
                render={({ field }) => (
                  <FormItem className="flex flex-col">
                    <FormLabel>Component</FormLabel>
                    <Popover open={popoverOpen} onOpenChange={setPopoverOpen} modal>
                      <PopoverTrigger asChild>
                        <FormControl>
                          <Button variant="outline" role="combobox" className={cn('justify-between', !field.value && 'text-muted-foreground')}>
                            {field.value ? declaredComponents?.find((component) => component.purl === field.value)?.name : 'Select component'}
                            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                          </Button>
                        </FormControl>
                      </PopoverTrigger>
                      <PopoverContent className="p-0">
                        <Command>
                          <CommandInput placeholder="Search component..." value={searchValue} onValueChange={setSearchValue} />
                          <CommandList>
                            <CommandEmpty>No components found.</CommandEmpty>
                            <CommandGroup>
                              <CommandItem asChild>
                                <div
                                  onClick={() => {
                                    setNewComponentDialogOpen(true);
                                    setPopoverOpen(false);
                                  }}
                                >
                                  <Plus className="mr-2 h-4 w-4" />
                                  Add new component
                                </div>
                              </CommandItem>
                              {searchValue ? (
                                <CommandItem asChild>
                                  <div onClick={handleSearchOnline} className="flex items-center justify-between gap-2">
                                    <div className="flex items-center gap-2">
                                      <Search className="h-3 w-3" />
                                      Search &ldquo;{searchValue}&rdquo; online
                                    </div>
                                    <ShortcutBadge shortcut="⌘ + Enter" />
                                  </div>
                                </CommandItem>
                              ) : null}
                            </CommandGroup>
                            <CommandSeparator />
                            <CommandGroup>
                              <ScrollArea className="h-52">
                                {declaredComponents?.map((component) => (
                                  <CommandItem
                                    key={component.purl}
                                    onSelect={() => {
                                      form.setValue('purl', component.purl);
                                      form.setValue('name', component.name);
                                      resetLicense();
                                      setPopoverOpen(false);
                                    }}
                                  >
                                    <Check className={cn('mr-2 h-4 w-4', component.purl === field.value ? 'opacity-100' : 'opacity-0')} />
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
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div>
                <Label>Purl</Label>
                <Input value={selectedPurl} disabled className="disabled:cursor-default" />
              </div>

              <FormItem className="flex flex-col">
                <Label>
                  License <span className="text-xs font-normal text-muted-foreground">(optional)</span>
                </Label>
                <SelectLicenseList key={licenseKey} onSelect={(value) => form.setValue('license', value)} />
              </FormItem>

              {withComment && (
                <FormField
                  control={form.control}
                  name="comment"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        Comment <span className="text-xs font-normal text-muted-foreground">(optional)</span>
                      </FormLabel>
                      <FormControl>
                        <Textarea {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              {filterBy === 'by_purl' && (
                <Alert>
                  <CircleAlert className="h-4 w-4" />
                  <AlertTitle>Be careful!</AlertTitle>
                  <AlertDescription>
                    <FilterByPurlList action={FilterAction.Replace} />
                  </AlertDescription>
                </Alert>
              )}

              <DialogFooter>
                <Button variant="ghost" onClick={onOpenChange}>
                  Cancel
                </Button>
                <Button ref={submitButtonRef} type="submit" onClick={form.handleSubmit(onSubmit)}>
                  Confirm <ShortcutBadge shortcut="⌘ + Enter" />
                </Button>
              </DialogFooter>
            </DialogContent>
          </form>
        </Dialog>
      </Form>
      {newComponentDialogOpen && (
        <NewComponentDialog onOpenChange={() => setNewComponentDialogOpen((prev) => !prev)} onCreated={onComponentSelected} />
      )}
      {onlineSearchDialogOpen && (
        <OnlineComponentSearchDialog
          onOpenChange={() => setOnlineSearchDialogOpen(false)}
          searchTerm={searchValue}
          onComponentSelect={async (c) => {
            onComponentSelected(c);
            await getMatchedLicenses(c.purl);
          }}
        />
      )}
    </>
  );
}
