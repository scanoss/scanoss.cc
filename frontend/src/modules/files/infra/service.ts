import files, { FileMatch } from '@/files';

export default class FileService {
  static getAll(): Record<string, FileMatch> {
    return files;
  }
}
