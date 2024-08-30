export interface GitFile {
  path: string;
}

export interface LocalFile {
  name: string;
  path: string;
  content: string;
  language: string | null;
}

export enum FilterAction {
  Ignore = 'ignore',
  Include = 'include',
  Remove = 'remove',
  Replace = 'replace',
}

export const filterActionLabelMap: Record<FilterAction, string> = {
  [FilterAction.Ignore]: 'Omit / Skip',
  [FilterAction.Include]: 'Include',
  [FilterAction.Remove]: 'Dismiss',
  [FilterAction.Replace]: 'Replace',
};
