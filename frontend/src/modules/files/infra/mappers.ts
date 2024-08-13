import { adapter } from 'wailsjs/go/models';

import { GitFile, LocalFile } from '../domain';

const languages: Record<string, string> = {
  sol: 'solidity',
  js: 'javascript',
  ts: 'typescript',
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

export const mapToGitFile = (response: adapter.GitFileDTO[]): GitFile[] =>
  response.map((file) => ({
    path: file.path,
  }));

export const mapToLocalFile = (response: adapter.FileDTO): LocalFile => {
  return {
    content: response.content,
    name: response.name,
    path: response.path,
    language: extractFileLanguage(response.path),
  };
};
