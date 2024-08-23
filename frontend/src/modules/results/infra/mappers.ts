import { common } from 'wailsjs/go/models';

import { MatchType, Result } from '../domain';

export const mapToResult = (response: common.ResultDTO[]): Result[] =>
  response.map((res) => ({
    path: res.file,
    matchType: res.matchType as MatchType,
    state: 'unstaged',
  }));
