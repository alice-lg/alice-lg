import { useMemo } from 'react';

import moment from 'moment';

import { useRouteServersMap }
  from 'app/context/route-servers';
import { useApiStatus }
  from 'app/context/api-status';
import { useRoutesLoading }
  from 'app/context/routes';
import { useSearchStatus }
  from 'app/context/search';

import RelativeTime
  from 'app/components/datetime/RelativeTime';


const RefreshIncomplete = () => {
  const routeServers = useRouteServersMap();
  const status = useApiStatus();
  const sources = status.store?.routes?.sources;

  let notInitialized = useMemo(() => {
    let missing = [];
    for (const id in sources) {
      if (sources[id].initialized) {
        continue;
      }
      missing.push(routeServers[id].name);
    }
    return missing;
  }, [routeServers, sources]);


  return (
    <>
    <p className="text-danger">
      Routes refresh was incomplete and results are missing
      from:
    </p>

      {notInitialized.map((name) => 
        <span key={name}>{name}<br /></span>
      )} 
      <br />
    </>
  );
}


const RefreshState = () => {
  const status = useApiStatus();
  if (!status.cachedAt || !status.ttlTime) {
    return null;
  }

  const cachedAt = moment.utc(status.cachedAt);
  const cacheTtl = moment.utc(status.ttlTime);

  const storeInitialized = status.store?.routes?.initialized === true;

  if (cacheTtl.isBefore(moment.utc())) {
    if (!storeInitialized) {
      return (
        <li>
          Routes cache is currently being refreshed. 
        </li>
      );
    }
    
    // This means cache is currently being rebuilt
    return (
      <li>
        Routes cache was built
          <b><RelativeTime
            fuzzyNow={5}
            pastEvent={true}
            value={cachedAt} /></b>
        and is currently being refreshed. 
      </li>
    );
  }

  if (!storeInitialized) {
    return (
      <li>
        <RefreshIncomplete />
        Next refresh in <b><RelativeTime value={cacheTtl} futureEvent={true} /></b>.
      </li>
    );
  }

  return (
    <li>
      Routes cache was built <b><RelativeTime fuzzyNow={5} value={cachedAt} /> </b>
      and will be refreshed <b><RelativeTime value={cacheTtl} futureEvent={true} /></b>.
    </li>
  );
}

const SearchStatus = () => {
  const isLoading = useRoutesLoading();
  const { queryDurationMs
        , totalReceived 
        , totalFiltered
        } = useSearchStatus();


  if (isLoading) {
    return null;
  }

  const queryDuration = queryDurationMs && queryDurationMs.toFixed(2);

  return (
    <div className="card">
      <div className="lookup-result-summary">
        <ul>
          <li>
            Found <b>{totalReceived}</b> received 
            and <b>{totalFiltered}</b> filtered routes.
          </li>
          <li>Query took <b>{queryDuration} ms</b> to complete.</li>
          <RefreshState />
        </ul>
      </div>
    </div>
  );
}

export default SearchStatus;
