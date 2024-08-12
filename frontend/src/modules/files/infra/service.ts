import {
  FileGetLocalContent,
  FileGetRemoteContent,
  GetFilesToBeCommited,
} from '../../../../wailsjs/go/main/App';
import { GitFile } from '../domain';
import { mapToGitFile, mapToLocalFile } from './mappers';

export default class FileService {
  static async getAllToBeCommited(): Promise<GitFile[]> {
    return GetFilesToBeCommited().then(mapToGitFile);
  }

  static async getLocalFileContent(path: string) {
    return FileGetLocalContent(path).then(mapToLocalFile);
  }

  static async getRemoteFileContent(path: string) {
    return FileGetRemoteContent(path).then(mapToLocalFile);
  }
}
