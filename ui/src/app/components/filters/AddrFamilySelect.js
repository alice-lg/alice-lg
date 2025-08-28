import { useAddrFamilyFilters } from 'app/context/filters';


/**
 * AddressFamilySelect filter IPv4 / IPv4 or keep both.
 */
const AddrFamilySelect = () => {
    const {filters, applyFilter, removeFilter} = useAddrFamilyFilters();
    const isApplied = (val) => filters.applied.find((f) => f.value === val);
    const isAvailable = (val) => filters.available.find((f) => f.value === val);

    let includeIp4 = !isApplied(2); // 2 = IPv6
    let includeIp6 = !isApplied(1); // 1 = IPv4
    let toggleIp4 = () => {
        if (includeIp4 && includeIp6) {
            applyFilter(2); // Apply IPv6 filter to hide IPv6
        } else {
            removeFilter(2);
        }
    };
    let toggleIp6 = () => {
        if (includeIp6 && includeIp4) {
            applyFilter(1); // Apply IPv4 filter to hide IPv4
        } else {
            removeFilter(1);
        }
    };

    if (!(isApplied(1) || isApplied(2))) {
        if (!(isAvailable(1) && isAvailable(2))) {
            return null;
        }
    }
    
    return (
        <div className="chk-filter-group">
          <label className="chk-filter-label">
            <input
              type="checkbox"
              className="chk-filter"
              checked={includeIp4}
              onChange={toggleIp4}  />
            IPv4
          </label>
          <label className="chk-filter-label">
            <input
              type="checkbox"
              className="chk-filter"
              checked={includeIp6}
              onChange={toggleIp6}  />
            IPv6
          </label>
        </div>
    );
};

export default AddrFamilySelect;
