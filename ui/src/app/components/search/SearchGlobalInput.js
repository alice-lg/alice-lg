import { useRef
       , useEffect
       }
  from 'react';

import { useQuery }
  from 'app/context/query';

import Content
  from 'app/components/content/Content';
import SearchQueryInput
  from 'app/components/search/SearchQueryInput';


/**
 * Help renders a quick help text
 */
const Help = () => {
  const [{q}] = useQuery();
  if (q) {
    return null;
  }
  return (
    <div className="lookup-help">
      <h3>Did you know?</h3>
      <p>You can search for</p>
      <ul>
        <li><b>Prefixes</b>,</li>
        <li><b>Peers</b> by entering their name and</li>
        <li><b>ASNs</b> by prefixing them with 'AS'</li>
      </ul>
      <p>Just start typing!</p>
    </div>
  );
}


/**
 * Global Search Input
 */
const SearchGlobalInput = () => {
  const ref = useRef();

  // Focus input
  useEffect(() => {
    if (ref.current) {
      ref.current.focus();
    }
  }, [ref]);

  return (
    <div className="lookup-container">
      <div className="card">
        <h2><Content id="lookup.title">Search on all route servers</Content></h2>
        <SearchQueryInput
          ref={ref}
          placeholder="Search for Prefixes, Peers or ASNs on all Route Servers" />
      </div>
      <Help /> 
    </div>
  );
}

export default SearchGlobalInput;
