import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { Check, ChevronsUpDown, Plus } from 'lucide-react';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { cn } from '@/lib/utils';
import useComponentFilterStore from '@/modules/components/stores/useComponentFilterStore';

import { GetDeclaredComponents } from '../../wailsjs/go/handlers/ComponentHandler';
import { entities } from '../../wailsjs/go/models';
import NewComponentDialog from './NewComponentDialog';
import { Button } from './ui/button';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from './ui/command';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from './ui/dialog';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from './ui/form';
import { Input } from './ui/input';
import { Label } from './ui/label';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { ScrollArea } from './ui/scroll-area';
import { useToast } from './ui/use-toast';

const ReplaceComponentFormSchema = z.object({
  purl: z.string().min(1, 'You must select a component.'),
});

interface ReplaceComponentDialogProps {
  onOpenChange: () => void;
  onReplaceComponent: () => void;
}

export default function ReplaceComponentDialog({
  onOpenChange,
  onReplaceComponent,
}: ReplaceComponentDialogProps) {
  const { toast } = useToast();
  const [newComponentDialogOpen, setNewComponentDialogOpen] = useState(false);
  const [declaredComponents, setDeclaredComponents] = useState<
    entities.DeclaredComponent[]
  >([]);

  const setReplaceWith = useComponentFilterStore(
    (state) => state.setReplaceWith
  );

  const form = useForm<z.infer<typeof ReplaceComponentFormSchema>>({
    resolver: zodResolver(ReplaceComponentFormSchema),
    defaultValues: {
      purl: '',
    },
  });

  const selectedPurl = form.watch('purl');

  const { data } = useQuery({
    queryKey: ['declaredComponents'],
    queryFn: GetDeclaredComponents,
  });

  const onSubmit = (values: z.infer<typeof ReplaceComponentFormSchema>) => {
    setReplaceWith(values.purl);

    setTimeout(() => {
      onReplaceComponent();
    }, 200);
  };

  const onComponentCreated = (component: entities.DeclaredComponent) => {
    const alreadyExists = declaredComponents.some(
      (c) => c.purl === component.purl
    );
    if (alreadyExists) {
      return toast({
        title: 'Warning',
        description: `A component with the same PURL already exists. Please select it from the list or enter a different PURL.`,
      });
    }
    setDeclaredComponents((prevState) => [...prevState, component]);
    form.setValue('purl', component.purl);
    setNewComponentDialogOpen(false);
  };

  useEffect(() => {
    if (data) {
      setDeclaredComponents(data);
    }
  }, [data]);

  return (
    <>
      <Form {...form}>
        <Dialog open onOpenChange={onOpenChange}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Replace</DialogTitle>
              <DialogDescription>
                You can search for an existing component or manually enter a
                PURL
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="purl"
                render={({ field }) => (
                  <FormItem className="flex flex-col">
                    <FormLabel>Component</FormLabel>
                    <Popover>
                      <PopoverTrigger asChild>
                        <FormControl>
                          <Button
                            variant="outline"
                            role="combobox"
                            className={cn(
                              'justify-between',
                              !field.value && 'text-muted-foreground'
                            )}
                          >
                            {field.value
                              ? declaredComponents?.find(
                                  (component) => component.purl === field.value
                                )?.name
                              : 'Select component'}
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
                                  onClick={() =>
                                    setNewComponentDialogOpen(true)
                                  }
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
                                    }}
                                  >
                                    <Check
                                      className={cn(
                                        'mr-2 h-4 w-4',
                                        component.purl === field.value
                                          ? 'opacity-100'
                                          : 'opacity-0'
                                      )}
                                    />
                                    <div>
                                      <p>{component.name}</p>
                                      <p className="text-xs text-muted-foreground">
                                        {component.purl}
                                      </p>
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
                <Input
                  value={selectedPurl}
                  disabled
                  className="disabled:cursor-default"
                />
              </div>

              <DialogFooter>
                <Button variant="ghost" onClick={onOpenChange}>
                  Cancel
                </Button>
                <Button type="submit">Confirm</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      </Form>
      {newComponentDialogOpen && (
        <NewComponentDialog
          onOpenChange={() => setNewComponentDialogOpen((prev) => !prev)}
          onCreated={onComponentCreated}
        />
      )}
    </>
  );
}
