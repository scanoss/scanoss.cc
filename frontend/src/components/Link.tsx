import { BrowserOpenURL } from '../../wailsjs/runtime';

export default function Link({ to }: { to: string }) {
  return (
    <a
      onClick={(e) => {
        e.preventDefault();
        BrowserOpenURL(to);
      }}
      className="cursor-pointer text-blue-500 underline"
    >
      {to}
    </a>
  );
}
