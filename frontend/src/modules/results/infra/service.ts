import { entities } from 'wailsjs/go/models';

import { ComponentGet, ResultGetAll } from '../../../../wailsjs/go/main/App';
import { MatchType, Result } from '../domain';
import { mapToResult } from './mappers';

export default class ResultService {
  static async getAll(matchType?: MatchType): Promise<Result[]> {
    return ResultGetAll({ matchType }).then(mapToResult);
  }

  static async getComponent(filePath: string): Promise<entities.ComponentDTO> {
    return ComponentGet(filePath);
  }
}
