import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { Check, ChevronsUpDown, CircleAlert, Plus } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import useKeyboardShortcut from '@/hooks/useKeyboardShortcut';
import { KEYBOARD_SHORTCUTS } from '@/lib/shortcuts';
import { cn } from '@/lib/utils';
import { FilterAction } from '@/modules/components/domain';
import { OnFilterComponentArgs } from '@/modules/components/stores/useComponentFilterStore';

import { entities } from '../../wailsjs/go/models';
import { GetDeclaredComponents } from '../../wailsjs/go/service/ComponentServiceImpl';
import FilterByPurlList from './FilterByPurlList';
import NewComponentDialog from './NewComponentDialog';
import SelectLicenseList from './SelectLicenseList';
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
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [newComponentDialogOpen, setNewComponentDialogOpen] = useState(false);
  const [declaredComponents, setDeclaredComponents] = useState<entities.DeclaredComponent[]>([]);
  const [licenseKey, setLicenseKey] = useState(0);

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

  const { data } = useQuery({
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

  const onComponentCreated = (component: entities.DeclaredComponent) => {
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

  useEffect(() => {
    if (data) {
      setDeclaredComponents(data);
    }
  }, [data]);

  const ref = useKeyboardShortcut(KEYBOARD_SHORTCUTS.confirm.keys, () => form.handleSubmit(onSubmit), {
    enableOnFormTags: true,
  });

  return (
    <>
      <Form {...form}>
        <Dialog open onOpenChange={onOpenChange}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <DialogContent ref={ref} tabIndex={-1}>
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
                    <Popover open={popoverOpen} onOpenChange={setPopoverOpen}>
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
                          <CommandInput placeholder="Search component..." />
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
                            </CommandGroup>
                            <CommandSeparator />
                            <CommandGroup>
                              <ScrollArea className="h-52">
                                {declaredComponents?.map((component) => (
                                  <CommandItem
                                    value={component.purl}
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
                <Label>License</Label>
                <SelectLicenseList key={licenseKey} onSelect={(value) => form.setValue('license', value)} />
              </FormItem>

              {withComment && (
                <FormField
                  control={form.control}
                  name="comment"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Comment</FormLabel>
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
                <Button type="submit" onClick={form.handleSubmit(onSubmit)}>
                  Confirm <span className="ml-2 rounded-sm bg-card p-1 text-[8px] leading-none">⌘ + Enter</span>
                </Button>
              </DialogFooter>
            </DialogContent>
          </form>
        </Dialog>
      </Form>
      {newComponentDialogOpen && (
        <NewComponentDialog onOpenChange={() => setNewComponentDialogOpen((prev) => !prev)} onCreated={onComponentCreated} />
      )}
    </>
  );
}
