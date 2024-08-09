import { ResultGetAll } from '../../../../wailsjs/go/main/App';
import { Result } from '../domain';
import { mapToResult } from './mappers';

export default class ResultService {
  static async getAll(): Promise<Result[]> {
    return ResultGetAll({}).then(mapToResult);
  }
}
