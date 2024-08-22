import { entities } from 'wailsjs/go/models';

export interface Result {
  path: string;
  matchType: MatchType;
}

export enum MatchType {
  File = 'file',
  Snippet = 'snippet',
}

export type Component = entities.ComponentDTO;
