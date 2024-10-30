import { useState } from 'react';

import useEnvironment from '@/hooks/useEnvironment';

import { Textarea } from './ui/textarea';

interface AddCommentTextareaProps {
  onChange: (value: string) => void;
}

export default function AddCommentTextarea({ onChange }: AddCommentTextareaProps) {
  const { modifierKey } = useEnvironment();
  const [inputValue, setInputValue] = useState('');

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputValue(e.target.value);
    onChange(e.target.value);
  };

  return (
    <div>
      <Textarea value={inputValue} onChange={handleChange} />
      <p className="mt-2 text-xs text-muted-foreground">
        Use <span className="rounded bg-primary px-1.5">{modifierKey.label} + return</span> to confirm
      </p>
    </div>
  );
}
