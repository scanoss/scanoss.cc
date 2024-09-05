import { entities } from 'wailsjs/go/models';

import { ComponentFilter } from '../../../../wailsjs/go/handlers/ComponentHandler';
import {
  FileGetLocalContent,
  FileGetRemoteContent,
} from '../../../../wailsjs/go/handlers/FileHandler';
import { SaveScanossBomFile } from '../../../../wailsjs/go/handlers/ScanossBomHandler';
import { mapToLocalFile } from './mappers';

export default class FileService {
  static async getLocalFileContent(path: string) {
    return FileGetLocalContent(path)
      .then(mapToLocalFile)
      .catch((e) => {
        throw new Error(e);
      });
  }

  static async getRemoteFileContent(path: string) {
    return FileGetRemoteContent(path)
      .then(mapToLocalFile)
      .catch((e) => {
        throw new Error(e);
      });
  }

  static async filterComponentByPath(dto: entities.ComponentFilterDTO) {
    return ComponentFilter(dto).catch((e) => {
      throw new Error(e);
    });
  }

  static async saveBomChanges() {
    return SaveScanossBomFile().catch((e) => {
      throw new Error(e);
    });
  }
}
