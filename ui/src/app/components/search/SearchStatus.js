import moment from 'moment';


import { useApiStatus }
  from 'app/context/api-status';
import { useRoutesLoading }
  from 'app/context/routes';
import { useSearchStatus }
  from 'app/context/search';

import RelativeTime
  from 'app/components/datetime/RelativeTime';




const RefreshState = () => {
  const status = useApiStatus();
  if (!status.cachedAt || !status.ttlTime) {
    return null;
  }

  const cachedAt = moment.utc(status.cachedAt);
  const cacheTtl = moment.utc(status.ttlTime);

  if (cacheTtl.isBefore(moment.utc())) {
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
