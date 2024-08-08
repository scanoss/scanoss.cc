export interface GitFile {
  path: string;
}

export interface LocalFile {
  name: string;
  path: string;
  content: string;
  language: string | null;
}
