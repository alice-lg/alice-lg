
import { useRef
       , useMemo
       }
  from 'react';

import { useNeighbors }
  from 'app/components/neighbors/Provider';
import NeighborsTable
  from 'app/components/neighbors/Table';

import LoadingIndicator
  from 'app/components/api/LoadingIndicator';

/**
 * Get AS from filter string
 */
const getFilterAsn = (filter) => {
  const tokens = filter.split("AS", 2);
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


const Neighbors = ({filter}) => {
  const refUp                  = useRef();
  const refDown                = useRef();
  const {isLoading, neighbors} = useNeighbors();

  const filtered = useMemo(
    () => filterNeighbors(neighbors, filter),
    [neighbors, filter]);

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

  // Group neighbors
  let neighborsUp = [];
  let neighborsDown = [];
  let neighborsIdle = [];

  for (let n of filtered) {
    let s = n.state.toLowerCase();
    if (s.includes("up") || s.includes("established") ) {
      neighborsUp.push(n);
    } else if (s.includes("down")) {
      neighborsDown.push(n);
    } else if (s.includes("start") || s.includes("active")) {
      neighborsIdle.push(n);
    } else {
      console.error("Couldn't classify neighbor by state:", n);
      neighborsUp.push(n);
    }
  }

  return (
    <>
      <div ref={refUp}>
        <NeighborsTable state="up"   neighbors={neighborsUp} />
      </div>
      <div ref={refDown}>
        <NeighborsTable state="idle" neighbors={neighborsIdle} />
        <NeighborsTable state="down" neighbors={neighborsDown} />
      </div>
    </>
  );
}

export default Neighbors;

