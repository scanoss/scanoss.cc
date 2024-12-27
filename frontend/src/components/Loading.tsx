interface Props {
  text?: string;
  size?: string;
  textSize?: string;
}

export default function Loading({ text = 'Loading', size = 'w-4 h-4', textSize = 'text-sm' }: Props) {
  return (
    <div className="inline-flex items-center justify-center leading-none">
      <svg className={`${size} animate-spin`} viewBox="0 0 16 16" xmlns="http://www.w3.org/2000/svg">
        <g className="opacity-40">
          <path
            fill="currentColor"
            d="M8 0a8 8 0 0 0-8 8 8 8 0 0 0 8 8 8 8 0 0 0 8-8 8 8 0 0 0-8-8zm0 14.4A6.4 6.4 0 0 1 1.6 8 6.4 6.4 0 0 1 8 1.6a6.4 6.4 0 0 1 6.4 6.4 6.4 6.4 0 0 1-6.4 6.4z"
          />
        </g>
        <path fill="currentColor" d="M8 0a8 8 0 0 1 8 8h-1.6A6.4 6.4 0 0 0 8 1.6V0z" />
      </svg>
      <span className={`ml-2 ${textSize} leading-none`}>{text}</span>
    </div>
  );
}
