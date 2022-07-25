
import { useMemo
       , useCallback
       }
  from 'react';

import { useAsnFilters }
  from 'app/context/filters';

import ButtonRemoveFilter
  from 'app/components/filters/ButtonRemoveFilter';


const PeersSelect = () => {
  const {filters, applyFilter, removeFilter} = useAsnFilters();
  const {applied, available} = filters;
  const active = applied[0]; // allow only one for now

  const sortedAvailable = useMemo(() =>
    available.sort((a, b) => (a.name.localeCompare(b.name))),
    [available]);

  const apply = useCallback((e) => {
    applyFilter(e.target.value);
  }, [applyFilter]);

  const removeActiveFilter = useCallback((filter) => {
    removeFilter(active.value);
  }, [removeFilter, active]);
    
  // Nothing to do if we don't have filters
  if (available.length === 0 && applied.length === 0) {
    return null;
  }

  if (active) {
    // Just render this, with a button for removal
    return (
      <table className="select-ctrl">
        <tbody>
          <tr>
            <td className="select-container">
              {active.name}
            </td>
            <td>
              <ButtonRemoveFilter onClick={removeActiveFilter} />
            </td>
          </tr>
        </tbody>
      </table>
    );
  }

  // Build options
  const optionsAvailable = sortedAvailable.map((filter) => {
    return (
      <option key={filter.value} value={filter.value}>
        {filter.name}, AS{filter.value} ({filter.cardinality})
      </option>
    );
  });

  return (
    <table className="select-ctrl">
      <tbody>
        <tr>
          <td className="select-container">
            <select className="form-control"
                    onChange={apply}
                    value={active?.value}>
              <option className="options-title"
                      value="none">Show only results from AS...</option>
              {optionsAvailable}
            </select>
          </td>
        </tr>
      </tbody>
    </table>
  );
}

export default PeersSelect;
