
import { useQueryParams }
  from 'app/context/query';
import { useRoutesReceived
       , useRoutesFiltered 
       , useRoutesNotExported 
       }
  from 'app/context/routes';


/**
 * Show a notice if no routes could be found
 */
const EmptyResults = () => {
  const query = useQueryParams({q: ""});

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
  const hasQuery = query.q !== "";

  if (isLoading) {
    return null; // We are not a loading indicator.
  }
 
  // Maybe this has something to do with a filter
  if (!hasContent && hasQuery) {
      return (
        <div className="card info-result-empty">
          <h4>No routes  matching your query.</h4>
          <p>Please check if your query is too restrictive.</p>
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
