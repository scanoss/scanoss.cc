import { entities } from 'wailsjs/go/models';

type ResultState = 'unstaged' | 'staged';
export interface Result {
  path: string;
  matchType: MatchType;
  state: ResultState;
}

export enum MatchType {
  File = 'file',
  Snippet = 'snippet',
}

export type Component = entities.ComponentDTO;
