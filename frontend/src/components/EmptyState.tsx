import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

interface EmptyStateProps {
  icon?: React.ReactNode;
  image?: string;
  title: string;
  subtitle?: string;
  action?: {
    label: string;
    onClick: () => void;
  };
  className?: string;
}

export default function EmptyState({ icon, image, title, subtitle, action, className }: EmptyStateProps) {
  return (
    <div className={cn('flex h-full flex-col items-center justify-center p-8', className)}>
      <div className="flex max-w-[420px] flex-col items-center text-center">
        {icon && <div className="mb-4 text-muted-foreground">{icon}</div>}
        {image && (
          <div className="mb-4">
            <img src={image} alt="" className="h-40 w-40 object-contain" />
          </div>
        )}
        <h3 className="mb-2 text-xl font-semibold">{title}</h3>
        {subtitle && <p className="mb-6 text-sm text-muted-foreground">{subtitle}</p>}
        {action && (
          <Button onClick={action.onClick} variant="default">
            {action.label}
          </Button>
        )}
      </div>
    </div>
  );
}
