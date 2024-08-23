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
