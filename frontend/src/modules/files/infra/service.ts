import { entities } from 'wailsjs/go/models';

import { ComponentFilter } from '../../../../wailsjs/go/handlers/ComponentHandler';
import {
  FileGetLocalContent,
  FileGetRemoteContent,
} from '../../../../wailsjs/go/handlers/FileHandler';
import { SaveScanossSettingsFile } from '../../../../wailsjs/go/handlers/ScanossSettingsHandler';
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
    return SaveScanossSettingsFile().catch((e) => {
      throw new Error(e);
    });
  }
}
