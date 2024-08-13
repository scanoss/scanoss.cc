import { adapter } from 'wailsjs/go/models';

export interface Result {
  path: string;
  matchType: MatchType;
}

export enum MatchType {
  File = 'file',
  Snippet = 'snippet',
  None = 'none',
}

export type Component = adapter.ComponentDTO;
