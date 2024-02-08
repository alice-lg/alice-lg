import { useMemo
       , useCallback
       }
  from 'react';

import { useQuery }
  from 'app/context/query';
import { useReadableCommunity }
  from 'app/context/bgp';
import { FILTER_GROUP_COMMUNITIES
       , FILTER_GROUP_EXT_COMMUNITIES
       , FILTER_GROUP_LARGE_COMMUNITIES
       , useFilters
       , useCommunitiesFilters
       , useExtCommunitiesFilters
       , useLargeCommunitiesFilters
       }
  from 'app/context/filters';

import ButtonRemoveFilter
  from 'app/components/filters/ButtonRemoveFilter';
import BgpCommunityLabel
  from 'app/components/routes/BgpCommunityLabel';


const AppliedCommunity = ({group, filter, onRemove}) => {
  const removeFilter = useCallback(() => {
    onRemove([group, filter.value]);
  }, [filter, group, onRemove]);
  const [{q}] = useQuery();
  let query = q || '';
  const repr = filter.value.join(':');
  const canRemove = !query.includes(repr);

  return (
    <tr>
      <td className="select-container">
        <BgpCommunityLabel community={filter.value} />
      </td>
      <td>
        {canRemove &&<ButtonRemoveFilter onClick={removeFilter} />}
      </td>
    </tr>
  );
}

const createCommunityOption = (group) => ({filter}) => {
  const name = useReadableCommunity(filter.value);

  if (filter.value.length === 0) {
    return null;
  }

  const cls = `select-bgp-community-0-${filter.value[0]} ` +
    `select-bgp-community-1-${filter.value[1]} ` +
    `select-bgp-community-2-${filter.value[2]}`;

  // Encode value to make stringly typed
  const value = JSON.stringify([group, filter.value]);
  return (
    <option value={value} className={cls}>
      {filter.name} {name} ({filter.cardinality})
    </option>
  );
}

const CommunityOption = createCommunityOption(
  FILTER_GROUP_COMMUNITIES);
const ExtCommunityOption = createCommunityOption(
  FILTER_GROUP_EXT_COMMUNITIES);
const LargeCommunityOption = createCommunityOption(
  FILTER_GROUP_LARGE_COMMUNITIES);


// Add / Remove a filter
const useUpdateFilters = (filter) => {
  const communities = useCommunitiesFilters();
  const extCommunities = useExtCommunitiesFilters();
  const largeCommunities = useLargeCommunitiesFilters();
  return useMemo(() => ({
    apply: {
      [FILTER_GROUP_COMMUNITIES]: communities.applyFilter,
      [FILTER_GROUP_EXT_COMMUNITIES]: extCommunities.applyFilter,
      [FILTER_GROUP_LARGE_COMMUNITIES]: largeCommunities.applyFilter,
    },
    remove: {
      [FILTER_GROUP_COMMUNITIES]: communities.removeFilter,
      [FILTER_GROUP_EXT_COMMUNITIES]: extCommunities.removeFilter,
      [FILTER_GROUP_LARGE_COMMUNITIES]: largeCommunities.removeFilter,
    }
  }), [communities, extCommunities, largeCommunities]);
}


const CommunitiesSelect = () => {
  const { filters } = useFilters();
  const { apply, remove } = useUpdateFilters();

  const communitiesFilters = useCommunitiesFilters();
  const extCommunitiesFilters = useExtCommunitiesFilters();
  const largeCommunitiesFilters = useLargeCommunitiesFilters();

  const applyFilter = useCallback((e) => {
    const [group, value] = JSON.parse(e.target.value); // filter
    apply[group](value);
  }, [apply]);

  const removeFilter = useCallback(([group, value]) => {
    remove[group](value);
  }, [remove]);

  // Nothing to do if we don't have filters or if the community
  // filter is disable because of a large result set.
  const filtersNotAvailable = filters.notAvailable;
  const isDisabled = filtersNotAvailable.includes("communities");

  const hasAvailable =
      communitiesFilters.filters.available.length > 0 ||
      extCommunitiesFilters.filters.available.length > 0 ||
      largeCommunitiesFilters.filters.available.length > 0;

  const communitiesAvailable =
    communitiesFilters.filters.available.sort((a, b) => {
      return (a.value[0] - b.value[0]) * 100000 + (a.value[1] - b.value[1]);
    });

  const extCommunitiesAvailable =
    extCommunitiesFilters.filters.available.sort((a, b) => {
      return (a.value[1] - b.value[1]) * 100000 + (a.value[2] - b.value[2]);
    });

  const largeCommunitiesAvailable = 
    largeCommunitiesFilters.filters.available.sort((a, b) => {
      return (a.value[0] - b.value[0]) * 10000000000 +
             (a.value[1] - b.value[1]) * 100000 +
             (a.value[2] - b.value[2]);
    });

  const communitiesOptions = communitiesAvailable.map((filter) => (
    <CommunityOption key={filter.value} filter={filter} />
  ));
  const extCommunitiesOptions = extCommunitiesAvailable.map((filter) => (
    <ExtCommunityOption key={filter.value} filter={filter} /> 
  ));
  const largeCommunitiesOptions = largeCommunitiesAvailable.map((filter) => (
    <LargeCommunityOption key={filter.value} filter={filter} />
  ));

  const appliedCommunities = communitiesFilters.filters.applied.map((filter) => (
    <AppliedCommunity 
      key={filter.value}
      group={FILTER_GROUP_COMMUNITIES}
      filter={filter}
      onRemove={removeFilter} />
  ));
  const appliedExtCommunities = extCommunitiesFilters.filters.applied.map((filter) => (
    <AppliedCommunity
      key={filter.value}
      group={FILTER_GROUP_EXT_COMMUNITIES}
      filter={filter}
      onRemove={removeFilter} />
  ));
  const appliedLargeCommunities = largeCommunitiesFilters.filters.applied.map((filter) => (
    <AppliedCommunity
      key={filter.value}
      group={FILTER_GROUP_LARGE_COMMUNITIES}
      filter={filter}
      onRemove={removeFilter}/>
  ));

  return (
    <table className="select-ctrl">
      <tbody>
        {appliedCommunities}
        {appliedExtCommunities}
        {appliedLargeCommunities}
        {isDisabled && <div className="text-muted">
            Due to a large number of results, selecting BGP communities
            becomes available only after selecting a route server or
            a neighbor.
          </div>}
        {hasAvailable &&
            <tr>
              <td className="select-container" colSpan="2">
                <select value="none"
                        onChange={applyFilter}
                        className="form-control">
                  <option value="none" className="options-title">
                    Select BGP Communities to match...
                  </option>
                  {communitiesOptions.length > 0 &&
                    <optgroup label="Communities">
                      {communitiesOptions}
                    </optgroup>}

                  {extCommunitiesOptions.length > 0 &&
                    <optgroup label="Ext. Communities">
                      {extCommunitiesOptions}
                    </optgroup>}

                  {largeCommunitiesOptions.length > 0 &&
                    <optgroup label="Large Communities">
                      {largeCommunitiesOptions}
                    </optgroup>}
                </select>
              </td>
            </tr>}
      </tbody>
    </table>
  );
}

export default CommunitiesSelect;
