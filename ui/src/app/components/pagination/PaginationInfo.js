
const PaginationInfo = ({results}) => {
  const totalResults = results.totalResults;
  const perPage = results.pageSize;
  const start = results.page * perPage + 1;
  const end = Math.min(start + perPage - 1, totalResults);

  if (results.totalPages <= 1) {
    let routes = "route";
    if (totalResults > 1) {
      routes = "routes";
    }
    return (
      <div className="routes-pagination-info pull-right">
        Showing <b>all</b> of <b>{totalResults}</b> {routes}
      </div>
    );
  }

  return (
    <div className="routes-pagination-info pull-right">
      Showing <b>{start} - {end}</b> of <b>{totalResults}</b> total routes
    </div>
   );
};

export default PaginationInfo;
