import { entities } from 'wailsjs/go/models';

import {ComponentFilter} from "../../../../wailsjs/go/handler/ComponentHandler";
import {FileGetLocalContent, FileGetRemoteContent} from "../../../../wailsjs/go/handler/FileHandler";

import { mapToLocalFile } from './mappers';
import {SaveScanossBomFile} from "../../../../wailsjs/go/handler/ScanossBomHandler";



export default class FileService {
/*  static async getAllToBeCommited(): Promise<GitFile[]> {
    return GetFilesToBeCommited()
      .then(mapToGitFile)
      .catch((e) => {
        throw new Error(e);
      });
  }*/

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
