import { LocalFile } from '@/modules/files/domain';

import { entities } from '../../../../wailsjs/go/models';
import FileDTO = entities.FileDTO;

const languages: Record<string, string> = {
  sol: 'solidity',
  js: 'javascript',
  ts: 'typescript',
  tsx: 'typescript',
  py: 'python',
  rb: 'ruby',
  sh: 'bash',
  go: 'go',
  java: 'java',
  c: 'c',
  cpp: 'cpp',
  h: 'c',
  hpp: 'cpp',
  cs: 'csharp',
  css: 'css',
  html: 'htmlbars',
  xml: 'xml',
  json: 'json',
  md: 'markdown',
  yml: 'yaml',
  scss: 'scss',
  less: 'less',
  sass: 'sass',
  sql: 'sql',
  txt: 'text',
};

const extractFileLanguage = (path: string): string | null => {
  const fileExtension = path.split('.').pop();

  if (!fileExtension) {
    return null;
  }

  return languages[fileExtension];
};

export const mapToLocalFile = (response: FileDTO): LocalFile => {
  return {
    content: response.content,
    name: response.name,
    path: response.path,
    language: extractFileLanguage(response.path),
  };
};
