/* eslint-disable @typescript-eslint/no-explicit-any */

interface WithErrorHandlingParams {
  asyncFn: (...args: any[]) => Promise<void>;
  onError: (error: unknown) => void;
  onFinish?: () => void;
}

export const withErrorHandling = ({ asyncFn, onError, onFinish }: WithErrorHandlingParams) => {
  return async (...args: any[]) => {
    try {
      await asyncFn(...args);
    } catch (error) {
      onError(error);
    } finally {
      onFinish?.();
    }
  };
};
