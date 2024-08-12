export interface Result {
  path: string;
  matchType: MatchType;
}

export enum MatchType {
  File = 'file',
  Snippet = 'snippet',
  None = 'none',
}
