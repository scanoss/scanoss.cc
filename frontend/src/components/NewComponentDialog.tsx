// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2024 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';

import { VALID_PURL_REGEX } from '@/modules/components/domain';

import { entities } from '../../wailsjs/go/models';
import { Button } from './ui/button';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from './ui/dialog';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from './ui/form';
import { Input } from './ui/input';

interface NewComponentDialogProps {
  onOpenChange: () => void;
  onCreated: (component: entities.DeclaredComponent) => void;
}

const NewComponentFormSchema = z.object({
  name: z.string().min(1, 'Required'),
  purl: z.string().regex(VALID_PURL_REGEX, 'Invalid PURL'),
});

export default function NewComponentDialog({ onOpenChange, onCreated }: NewComponentDialogProps) {
  const form = useForm<z.infer<typeof NewComponentFormSchema>>({
    resolver: zodResolver(NewComponentFormSchema),
    defaultValues: {
      name: '',
      purl: '',
    },
  });

  const onSubmit = (values: z.infer<typeof NewComponentFormSchema>) => {
    onCreated(values);
  };

  return (
    <Form {...form}>
      <Dialog open onOpenChange={onOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create Component</DialogTitle>
          </DialogHeader>

          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="purl"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>PURL</FormLabel>
                  <FormControl>
                    <Input {...field} type="text" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button variant="ghost" onClick={onOpenChange}>
                Cancel
              </Button>
              <Button type="submit">Create</Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </Form>
  );
}
