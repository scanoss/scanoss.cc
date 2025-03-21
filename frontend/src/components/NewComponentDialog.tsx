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
        <DialogContent className="p-4">
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
