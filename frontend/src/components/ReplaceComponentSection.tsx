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
import { useEffect, useState } from 'react';

import { entities } from '../../wailsjs/go/models';
import { GetDeclaredComponents } from '../../wailsjs/go/service/ComponentServiceImpl';
import { GetLicensesByPurl } from '../../wailsjs/go/service/LicenseServiceImpl';
import ComponentSelector from './ComponentSelector';
import NewComponentDialog from './NewComponentDialog';
import OnlineComponentSearchDialog from './OnlineComponentSearchDialog';
import SelectLicenseList from './SelectLicenseList';
import { Label } from './ui/label';
import { useToast } from './ui/use-toast';

interface ReplaceComponentSectionProps {
  licenseKey: number;
  onComponentChange: (component: entities.DeclaredComponent | null) => void;
  onLicenseChange: (license: string | undefined) => void;
}

export default function ReplaceComponentSection({
  licenseKey,
  onComponentChange,
  onLicenseChange,
}: ReplaceComponentSectionProps) {
  const { toast } = useToast();

  const [newComponentDialogOpen, setNewComponentDialogOpen] = useState(false);
  const [onlineSearchDialogOpen, setOnlineSearchDialogOpen] = useState(false);
  const [onlineSearchTerm, setOnlineSearchTerm] = useState('');
  const [declaredComponents, setDeclaredComponents] = useState<entities.DeclaredComponent[]>([]);
  const [selectedComponent, setSelectedComponent] = useState<entities.DeclaredComponent | null>(null);
  const [matchedLicenses, setMatchedLicenses] = useState<entities.License[]>([]);
  const [internalLicenseKey, setInternalLicenseKey] = useState(0);

  const { data: componentsData, error: componentsError } = useQuery({
    queryKey: ['declaredComponents'],
    queryFn: GetDeclaredComponents,
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

  const handleComponentSelected = async (component: entities.DeclaredComponent) => {
    const alreadyExists = declaredComponents.some((c) => c.purl === component.purl);
    if (!alreadyExists) {
      setDeclaredComponents((prev) => [...prev, component]);
    }
    setSelectedComponent(component);
    onComponentChange(component);
    onLicenseChange(undefined);
    setInternalLicenseKey((k) => k + 1);
    setNewComponentDialogOpen(false);
    await getMatchedLicenses(component.purl);
  };

  // Reset when licenseKey changes (modal reopened)
  useEffect(() => {
    setSelectedComponent(null);
    setMatchedLicenses([]);
  }, [licenseKey]);

  return (
    <>
      <div className="flex flex-col gap-2">
        <Label>Replace with</Label>
        <ComponentSelector
          components={declaredComponents}
          selectedComponent={selectedComponent}
          onSelect={handleComponentSelected}
          onAddNew={() => setNewComponentDialogOpen(true)}
          onSearchOnline={(term) => {
            setOnlineSearchTerm(term);
            setOnlineSearchDialogOpen(true);
          }}
        />
        {selectedComponent && <p className="font-mono text-xs text-muted-foreground">{selectedComponent.purl}</p>}

        <div className="mt-2 flex flex-col gap-2">
          <Label>
            License <span className="text-xs font-normal text-muted-foreground">(optional)</span>
          </Label>
          <SelectLicenseList key={`${licenseKey}-${internalLicenseKey}`} onSelect={onLicenseChange} matchedLicenses={matchedLicenses} />
        </div>
      </div>

      {newComponentDialogOpen && (
        <NewComponentDialog
          onOpenChange={() => setNewComponentDialogOpen(false)}
          onCreated={handleComponentSelected}
        />
      )}
      {onlineSearchDialogOpen && (
        <OnlineComponentSearchDialog
          onOpenChange={() => setOnlineSearchDialogOpen(false)}
          searchTerm={onlineSearchTerm}
          onComponentSelect={handleComponentSelected}
        />
      )}
    </>
  );
}
