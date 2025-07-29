import { ExternalLink } from 'lucide-react';

import { entities } from '../../wailsjs/go/models';
import Link from './Link';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './ui/table';

export function ComponentSearchTable({
  components,
  onComponentSelect,
}: {
  components: entities.SearchedComponent[];
  onComponentSelect: (component: entities.SearchedComponent) => void;
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Component</TableHead>
          <TableHead>Purl</TableHead>
          <TableHead>URL</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {components.map((component) => (
          <TableRow key={component.purl} onClick={() => onComponentSelect(component)} className="cursor-pointer">
            <TableCell className="font-medium">{component.component}</TableCell>
            <TableCell>{component.purl}</TableCell>
            <TableCell>
              <Link to={component.url}>
                <div className="inline-flex items-center">
                  {component.url}
                  <ExternalLink className="ml-1 h-3 w-3" />
                </div>
              </Link>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
