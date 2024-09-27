import { entities } from 'wailsjs/go/models';

import { ComponentGet } from '../../../../wailsjs/go/handlers/ComponentHandler';
import { ResultGetAll } from '../../../../wailsjs/go/handlers/ResultHandler';
import { MatchType } from '../domain';

export default class ResultService {
  static async getAll(
    matchType?: MatchType,
    query?: string
  ): Promise<entities.ResultDTO[]> {
    return ResultGetAll({ match_type: matchType, query }).catch((e) => {
      throw new Error(e);
    });
  }

  static async getComponent(filePath: string): Promise<entities.ComponentDTO> {
    return ComponentGet(filePath).catch((e) => {
      throw new Error(e);
    });
  }
}
