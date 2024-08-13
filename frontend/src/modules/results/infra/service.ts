import { ResultGetAll } from '../../../../wailsjs/go/main/App';
import { MatchType, Result } from '../domain';
import { mapToResult } from './mappers';

export default class ResultService {
  static async getAll(matchType?: MatchType): Promise<Result[]> {
    return ResultGetAll({ matchType }).then(mapToResult);
  }
}
