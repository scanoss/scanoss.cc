export enum FilterAction {
  Ignore = 'ignore',
  Include = 'include',
  Remove = 'remove',
  Replace = 'replace',
}

export type FilterBy = 'path' | 'purl';

export const filterActionLabelMap: Record<FilterAction, string> = {
  [FilterAction.Ignore]: 'Omit / Skip',
  [FilterAction.Include]: 'Include',
  [FilterAction.Remove]: 'Dismiss',
  [FilterAction.Replace]: 'Replace',
};

export interface OnAddFilterArgs {
  action: FilterAction;
  filterBy: 'by_file' | 'by_purl';
  withComment?: boolean;
}
