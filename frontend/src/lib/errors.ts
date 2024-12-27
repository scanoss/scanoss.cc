/* eslint-disable @typescript-eslint/no-explicit-any */

interface WithErrorHandlingParams {
  asyncFn: (...args: any[]) => Promise<void>;
  onError: (error: unknown) => void;
  onFinish?: () => void;
  onSuccess?: () => Promise<void>;
}

export const withErrorHandling = ({ asyncFn, onError, onFinish, onSuccess }: WithErrorHandlingParams) => {
  return async (...args: any[]) => {
    try {
      await asyncFn(...args);
      await onSuccess?.();
    } catch (error) {
      onError(error);
    } finally {
      onFinish?.();
    }
  };
};
