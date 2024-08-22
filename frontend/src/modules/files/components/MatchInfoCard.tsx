import { Component } from '@/modules/results/domain';

interface MatchInfoCardProps {
  component: Component;
}

export default function MatchInfoCard({ component }: MatchInfoCardProps) {
  return (
    <div className="flex cursor-pointer items-center justify-between rounded-sm border border-primary-foreground bg-primary p-3">
      <div className="flex items-center gap-8">
        <div>
          <div className="font-bold text-primary-foreground">
            {component.component}
          </div>
          <div className="text-sm">{component.purl?.[0]}</div>
        </div>
        {component.version && (
          <div>
            <div className="text-sm text-muted-foreground">Version</div>
            <div className="text-sm">{component.version}</div>
          </div>
        )}
        {component.licenses?.length ? (
          <div>
            <div className="text-sm text-muted-foreground">License</div>
            <div className="text-sm">{component.licenses?.[0].name}</div>
          </div>
        ) : null}
        <div>
          <div className="text-sm text-muted-foreground">Detected</div>
          <div className="text-sm text-green-500 first-letter:uppercase">
            {component.id}
          </div>
        </div>
        <div>
          <div className="text-sm text-muted-foreground">Match</div>
          <div className="text-sm">{component.matched}</div>
        </div>
      </div>
    </div>
  );
}
