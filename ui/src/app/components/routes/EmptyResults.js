
import { useQuery }
  from 'app/context/query';
import { useRoutesReceived
       , useRoutesFiltered 
       , useRoutesNotExported 
       }
  from 'app/context/routes';
import { isTimeoutError }
  from 'app/context/errors';


/**
 * Show an error if present
 */
const ErrorResult = ({error}) => {
  const info = error.response?.data;
  if (!info) {
    return null;
  }
  
  return (
    <p className="text-danger">Reason: {info.message}</p>
  );
}


/**
 * Show a notice if no routes could be found
 */
const EmptyResults = () => {
  const [query] = useQuery({q: ""});

  const received = useRoutesReceived();
  const filtered = useRoutesFiltered();
  const notExported = useRoutesNotExported();

  // Conditions
  const hasContent = received.totalResults > 0 ||
                     filtered.totalResults > 0 ||
                     notExported.totalResults > 0;
  const isLoading = received.loading ||
                    filtered.loading ||
                    notExported.loading;
  const isRequested = received.requested ||
                      filtered.requested ||
                      notExported.requested;
  const hasQuery = query.q !== "";

  if (isLoading) {
    return null; // We are not a loading indicator.
  }
 
  // Maybe this has something to do with a filter
  if (!hasContent && hasQuery && isRequested) {
      if (isTimeoutError(received?.error)) {
        return (
          <div className="card info-result-empty">
            <h4 className="text-danger">The query took too long to process.</h4>
            <p>
              Unfortunately, it looks like the query matches a lot of routes.<br />
              Please try to refine your query to be more specific.
            </p>
          </div>
        );
      }

      return (
        <div className="card info-result-empty">
          <h4>No routes  matching your query.</h4>
          <p>Please check if your query is too restrictive.</p>
          {received?.error && <ErrorResult error={received.error} />}
        </div>
      );
  }

  if (hasContent) {
    return null; // Nothing to do then.
  }

  return (
    <div className="card info-result-empty">
      <p className="card-body">
          There are <b>no routes</b> to display for this neighbor.
      </p>
    </div>
  );
}

export default EmptyResults;
