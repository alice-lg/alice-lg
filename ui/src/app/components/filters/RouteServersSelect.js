
import { useMemo
       , useCallback
       }
  from 'react';

import { useSourceFilters }
  from 'app/context/filters';


const RouteServersSelect = () => {
  const {filters, applyFilter} = useSourceFilters();
  const {applied, available} = filters;
  const active = applied[0];

  const sortedAvailable = useMemo(() =>
    available.sort((a, b) => (a.value - b.value)),
    [available]);

  const selectFilter = useCallback((e) => {
    applyFilter(e.target.value);
  }, [applyFilter]);

  const removeFilter = useCallback((filter) => {
    console.log("removing filter:", filter);
  }, []);
    
  // Nothing to do if we don't have filters
  if (available.length === 0) {
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
              <button className="btn btn-remove"
                      onClick={removeFilter}>
                <i className="fa fa-times" />
              </button>
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
        {filter.name} ({filter.cardinality})
      </option>
    );
  });

  return (
    <table className="select-ctrl">
      <tbody>
        <tr>
          <td className="select-container">
            <select className="form-control"
                    onChange={selectFilter}
                    value={active.value}>
              <option value="none" className="options-title">Show results from RS...</option>
              {optionsAvailable}
            </select>
          </td>
        </tr>
      </tbody>
    </table>
  );
};

export default RouteServersSelect;
