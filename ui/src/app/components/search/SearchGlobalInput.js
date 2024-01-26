import { useRef
       , useEffect
       }
  from 'react';

import { useQuery }
  from 'app/context/query';
import { useContent }
  from 'app/context/content';

import Content
  from 'app/components/content/Content';
import SearchQueryInput
  from 'app/components/search/SearchQueryInput';
import BgpCommunityLabel
  from 'app/components/routes/BgpCommunityLabel';

/**
 * Lookup Example
 */
const LookupExample = ({example}) => {
  const type = example[0];
  const value = example.slice(1);

  const communityURL = (value) =>
    "/search?q=" + encodeURIComponent(`#${value.join(":")}`);

  switch (type) {
    case "community":
      return (
        <li className="community">
          <a href={communityURL(value)}>
            <BgpCommunityLabel community={value} />
          </a>
        </li>
      );
    default:
      return (
        <li className={type}>
          <a href={`/search?q=${value}`}>
            <span className={`label label-default label-${type}`}>{value}</span>
          </a>
        </li>
      );
  };
}

/**
 * Lookup Examples
 */
const LookupExamples = () => {
  const content = useContent();

  let examples = content.lookup?.examples;
  if (!examples) {
    return null;
  }

  return (
    <div className="lookup-examples">
      <h3>Some Examples</h3>
      <ul>
        {examples.map((example, i) =>
          <LookupExample key={i} example={example} />)}
      </ul>
    </div>
  )
}


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
        <li><b>Communities</b> by prefixing them with '#'</li>
      </ul>
      <p>Just start typing!</p>
      <LookupExamples />
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
