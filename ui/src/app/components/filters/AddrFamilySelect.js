import { useAddrFamilyFilters } from 'app/context/filters';


/**
 * AddressFamilySelect filter IPv4 / IPv4 or keep both.
 */
const AddrFamilySelect = () => {
    const {filters, applyFilter, removeFilter} = useAddrFamilyFilters();

    const isApplied = (val) => filters.applied.find((f) => f.value === val);
  
    let showIp4 = !isApplied(2); // 2 = IPv6
    let showIp6 = !isApplied(1); // 1 = IPv4
    let toggleIp4 = () => {
        if (showIp4 && showIp6) {
            applyFilter(2); // Apply IPv6 filter to hide IPv6
        } else {
            removeFilter(2);
        }
    };
    let toggleIp6 = () => {
        if (showIp6 && showIp4) {
            applyFilter(1); // Apply IPv4 filter to hide IPv4
        } else {
            removeFilter(1);
        }
    };
    

    return (
        <div>
          <label className="chk-filter-label">
            <input
              type="checkbox"
              className="chk-filter"
              checked={showIp4}
              onChange={toggleIp4}  />
            IPv4
          </label>
          <label>
            <input
              type="checkbox"
              className="chk-filter"
              checked={showIp6}
              onChange={toggleIp6}  />
            IPv6
          </label>
        </div>
    );
};

export default AddrFamilySelect;
