
import { useLocation }
  from 'react-router-dom';

import { useRef
       , useMemo
       , useEffect
       }
  from 'react';

import { useQuery }
  from 'app/context/query';
import { useNeighbors }
  from 'app/context/neighbors';

import NeighborsTable
  from 'app/components/neighbors/NeighborsTable';
import LoadingIndicator
  from 'app/components/spinners/LoadingIndicator';

/**
 * Get AS from filter string
 */
const getFilterAsn = (filter) => {
  filter = filter.toLowerCase();
  const tokens = filter.split("as", 2);
  if (tokens.length !== 2) {
    return false; // Not an ASN query
  }
  const asn = parseInt(tokens[1], 10);
  // Check if ASN is a valid number
  if (!(asn >= 0)) {
    return false;
  }
  return asn;
}


/**
 * Filter neighbors 
 */
const filterNeighbors = (protocols, filter) => {
  let filtered = [];
  if (!filter || filter === "") {
    return protocols; // nothing to do here
  }

  // We support different filter modes:
  // - Default: Try to match as much as possible
  // - AS$num: Try to match ASN only
  const filterAsn = getFilterAsn(filter);
  if (filterAsn) {
    filtered = protocols.filter((p) => {
      return (p.asn === filterAsn);
    });
  } else {
    filter = filter.toLowerCase();
    filtered = protocols.filter((p) => {
      return (p.asn === filter ||
              p.address.toLowerCase().indexOf(filter) !== -1 ||
              p.description.toLowerCase().indexOf(filter) !== -1);
    });
  }

  return filtered;
}


const Neighbors = () => {
  const { hash }  = useLocation();
  const [ query ] = useQuery();
  const filter    = query.q;

  const refUp                  = useRef();
  const refDown                = useRef();

  const {isLoading, neighbors} = useNeighbors();

  const filtered = useMemo(
    () => filterNeighbors(neighbors, filter),
    [neighbors, filter]);

  // Group neighbors
  const groups = useMemo(() => {
    let up = [];
    let down = [];
    let idle = [];

    for (let n of filtered) {
      let s = n.state.toLowerCase();
      if (s.includes("up") || s.includes("established") ) {
        up.push(n);
      } else if (s.includes("down")) {
        down.push(n);
      } else if (s.includes("start") || s.includes("active")) {
        idle.push(n);
      } else {
        console.error("Couldn't classify neighbor by state:", n);
        up.push(n);
      }
    }
    return {up, down, idle};
  }, [filtered]);

  // Scroll to anchor
  useEffect(() => {
    if (hash === "#sessions-down") {
      if (!refDown.current) {
        return;
      }
      refDown.current.scrollIntoView();
    }
    if (hash === "#sessions-up") {
      if (!refUp.current) {
        return;
      }
      refUp.current.scrollIntoView();
    }
  }, [hash, refDown, refUp, filtered]);

  if (isLoading) {
    return <LoadingIndicator show={true} />;
  }

  if (!filtered || filtered.length === 0) {
    // Empty Neighbors List
    return (
      <div className="card">
        <p className="help-block">
          No neighbors could be found.
        </p>
      </div>
    );
  }

  return (
    <>
      <div ref={refUp}>
        <NeighborsTable state="up"   neighbors={groups.up} />
      </div>
      <div ref={refDown}>
        <NeighborsTable state="idle" neighbors={groups.idle} />
        <NeighborsTable state="down" neighbors={groups.down} />
      </div>
    </>
  );
}

export default Neighbors;

