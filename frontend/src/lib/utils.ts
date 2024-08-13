import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function encodeFilePath(filePath: string) {
  return btoa(filePath);
}

export function decodeFilePath(encodedPath: string) {
  return atob(encodedPath);
}
