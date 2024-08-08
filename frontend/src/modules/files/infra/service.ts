import {
  GetFilesToBeCommited,
  GetLocalFileContent,
} from '../../../../wailsjs/go/main/App';
import { GitFile } from '../domain';
import { mapToGitFile, mapToLocalFile } from './mappers';

export default class FileService {
  static async getAllToBeCommited(): Promise<GitFile[]> {
    return GetFilesToBeCommited().then(mapToGitFile);
  }

  static async getLocalFileContent(path: string) {
    return GetLocalFileContent(path).then(mapToLocalFile);
  }
}
