
import { useMemo }
  from 'react';

import { Link
       , useLocation
       } 
  from 'react-router-dom';

import { useRoutesReceived
       , useRoutesFiltered
       , useRoutesNotExported
       }
  from 'app/context/routes';


const shouldShowLink = ({ loading, totalResults }) => {
  return (!loading && totalResults > 0);
}


/*
 * Quick links:
 * Jump to anchors for: not exported, filtered and received
 */
const QuickLinks = () => {
  const location = useLocation();

  const received = useRoutesReceived();
  const filtered = useRoutesFiltered();
  const notExported = useRoutesNotExported();

  const locRecevied = useMemo(() => ({...location, hash: "routes-received"}), [
    location,
  ]);
  const locFiltered = useMemo(() => ({...location, hash: "routes-filtered"}), [
    location,
  ]);
  const locNotExported = useMemo(() => ({...location, hash: "routes-not-exported"}), [
    location,
  ]);

  const showReceived = shouldShowLink(received);
  const showFiltered = shouldShowLink(filtered);
  const showNotExported = shouldShowLink(notExported);

  if (!showReceived && !showFiltered && !showNotExported) {
    return null;
  }

  return (
    <div className="quick-links routes-quick-links">
      <span>Go to:</span>
      <ul>
        { showFiltered &&
          <li className="filtered">
            <Link to={locFiltered}>Filtered</Link></li>}
        { showReceived &&
          <li className="received">
            <Link to={locRecevied}>Accepted</Link></li>}
        { showNotExported &&
          <li className="not-exported">
            <Link to={locNotExported}>Not Exported</Link></li>}
      </ul>
    </div>
  );
}

export default QuickLinks;
