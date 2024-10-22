import ActionToolbar from './ActionToolbar';
import FilterComponentActions from './FilterComponentActions';

export default function Header() {
  return (
    <div className="grid h-full grid-cols-2 md:grid-cols-3">
      <div className="hidden md:block"></div>
      <FilterComponentActions />
      <ActionToolbar />
    </div>
  );
}
