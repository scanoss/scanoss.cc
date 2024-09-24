import { entities } from 'wailsjs/go/models';

import { MatchType, Result } from '../domain';

export const mapToResult = (response: entities.ResultDTO[]): Result[] =>
  response.map((res) => ({
    path: res.path,
    matchType: res.match_type as MatchType,
    state: 'unstaged',
    bomState: null,
  }));
