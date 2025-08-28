import { useTotalFilters }
  from 'app/context/filters';

import RouteServersSelect
  from 'app/components/filters/RouteServersSelect';
import PeersSelect
  from 'app/components/filters/PeersSelect';
import CommunitiesSelect
  from 'app/components/filters/CommunitiesSelect';
import AddrFamilySelect
  from 'app/components/filters/AddrFamilySelect';


const withGroup = (title, FilterGroup) => (props) => {
  const content = FilterGroup(props); 
  if (content === null) {
    return null;
  }
  return (
    <div className="filter-editor-widget">
      <h2>{title}</h2>
      {content}
    </div>
  );
}

const RouteServersGroup = withGroup("Route Server", RouteServersSelect);
const CommunitiesGroup = withGroup("BGP Communities", CommunitiesSelect);
const PeersGroup = withGroup("Neighbor", PeersSelect);
const AddrFamilyGroup = withGroup("Address Family", AddrFamilySelect);

const FiltersEditor = () => {
  const totalFilters = useTotalFilters();
  if (totalFilters === 0) {
    return null;
  }

  return (
    <div className="card lookup-filters-editor">
      <AddrFamilyGroup />
      <RouteServersGroup />
      <PeersGroup />
      <CommunitiesGroup />
    </div>
  );
}

export default FiltersEditor;
