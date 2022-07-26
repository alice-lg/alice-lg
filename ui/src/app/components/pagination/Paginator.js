
import { useMemo
       , useCallback
       }
  from 'react';

import { Link
       , useNavigate 
       }
  from 'react-router-dom';

import { useMakePageLocation 
       , usePageQuery
       }
  from 'app/context/pagination';


const makePages = (total, max) => {
  const pages = Array.from(Array(total), (_, i) => i);
  return {
    items:  pages.slice(0, max),
    select: pages.slice(max),
  }
}

const PageLink = ({
  page,
  to,
  disabled = false,
  label = false,
}) => {
  const linkLabel = label || `${page + 1}`;
  if (disabled) {
    return <span>{linkLabel}</span>;
  }

  return (
    <Link to={to}>{linkLabel}</Link>
  );
}

/**
 * Render a drop down select
 */
const PageSelect = ({page, pages, onSelect}) => {
  const handleChange = useCallback((e) => onSelect(e.target.value), [
    onSelect,
  ]);

  if (pages.length === 0) {
    return null; // nothing to do here.
  }

  const items = pages.map((p) => (
    <option key={p} value={p}>{p + 1}</option>
  ));

  const active = page >= pages[0];
  let itemClassName = "";
  if (active) {
    itemClassName = "active";
  }

  return (
    <li className={itemClassName}>
      <select className="form-control pagination-select"
              value={page}
              onChange={handleChange}>
        { page < pages[0] && <option value={pages[0]}>more...</option> }
        {items}
      </select>
    </li>
  );
}

const Paginator = ({
  results,
  pageKey,
  maxItems=12,
  anchor="",
}) => {
  const navigate = useNavigate();
  const query = usePageQuery();
  const makePageLocation = useMakePageLocation(pageKey, anchor);
  const current = query[pageKey];

  const pages = useMemo(() => makePages(results.totalPages, maxItems), [
    results, maxItems,
  ]);

  // Callback: Page Selected
  const selectPage = useCallback((p) => {
    navigate(makePageLocation(p));
  }, [navigate, makePageLocation]);

  // Render list of page items
  const pageLinks = pages.items.map((p) => {
    const to = makePageLocation(p);
    let className = "";
    if (current === p) {
      className = "active";
    }
    return (
      <li key={p} className={className}>
        <PageLink page={p} to={to} />
      </li>
    );
  });

  // Links classes
  let prevLinkClass = "";
  if (current === 0) {
    prevLinkClass = "disabled";
  }

  let nextLinkClass = "";
  if (current + 1 === results.totalPages) {
    nextLinkClass = "disabled";
  }

  const toPrevious = makePageLocation(current - 1);
  const toNext = makePageLocation(current + 1);

  if (results.totalPages <= 1) {
    return null; // Nothing to paginate
  }

  return (
    <nav aria-label="Routes Pagination">
      <ul className="pagination">
        <li className={prevLinkClass}>
          <PageLink
            to={toPrevious}
            page={current - 1}
            label="&laquo;"
            disabled={current === 0} />
        </li>

        {pageLinks}

        <PageSelect pages={pages.select}
                    page={current}
                    onSelect={selectPage} />

        {pages.select.length === 0 &&
          <li className={nextLinkClass}>
            <PageLink 
              to={toNext}
              page={current + 1}
              disabled={(current + 1) === results.totalPages}
              label="&raquo;" />
          </li>}
      </ul>
    </nav>
  );
}

export default Paginator;
