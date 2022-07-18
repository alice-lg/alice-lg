
import bigInt from 'big-integer';

import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faCircleArrowUp
       , faCircleArrowDown
       }
  from '@fortawesome/free-solid-svg-icons';


import { useMemo }
  from 'react';

import { useParams
       , Link
       }
  from 'react-router-dom';

import { ipToNumeric }
  from 'app/utils/ip'

import { useConfig }
  from 'app/context/config';
import { useRouteServer }
  from 'app/context/route-servers';
import { useQuery
       , useQueryLink
       }
  from 'app/context/query';

import RelativeTimestamp
  from 'app/components/datetime/RelativeTimestamp';

/**
 * Default: Sort by ASN, ascending order.
 */
const querySortDefault = {
  s: 'asn',
  o: 'asc',
};

const lookupProperty = (obj, path) => {
  let property = path.split(".").reduce((acc, part) => acc[part], obj);
  if (typeof(property) == "undefined") {
    property = `Property "${path}" not found in object.`;
  }

  return property;
}

/**
 * Check if state is up or established
 */
const isUpState = (s) => {
    if (!s) { return false; }
    s = s.toLowerCase();
    return (s.includes("up") || s.includes("established"));
}

/**
 * Sort alphanumeric
 */
const sortAnum = (sort) => {
  return (a, b) => {
    const va = a[sort];
    const vb = b[sort];
    if (va < vb ) { return -1; }
    if (va > vb ) { return 1;  }
    return 0;
  }
}

/**
 * Sort by IPAddress
 */
const sortIpAddr = (sort) => {
  return (a, b) => {
    const va = ipToNumeric(a[sort]);
    const vb = ipToNumeric(b[sort]);

    // Handle ipv6 case
    if (va instanceof bigInt) {
      return va.compareTo(vb);
    }

    if (va < vb ) { return -1; }
    if (va > vb ) { return 1;  }
    return 0;
  }
}

/**
 * Sort with order
 */
const sortOrder = (cmp, order) => {
  return (a, b) => {
    const res = cmp(a, b);
    if (order === 'desc') {
      return res * -1;
    }
    return res;
  }
}

/**
 * Sort neighbors
 */
const sortNeighbors = (neighbors, sort, order) => {
  // Make compare function
  let cmp = sortAnum(sort);
  if (sort === "address") {
    cmp = sortIpAddr(sort);
  }
  return neighbors.sort(sortOrder(cmp, order));
}


/**
 * Section renders the sections title
 */
const Section = ({state}) => {
  let sectionTitle = '';
  let sectionCls = 'card-header card-header-neighbors ';
  switch(state) {
    case 'up':
      sectionTitle  = 'BGP Sessions Established';
      sectionCls  += 'established ';
      break;
    case 'down':
      sectionTitle = 'BGP Sessions Down';
      sectionCls += 'down ';
      break;
    case 'start':
      sectionTitle = 'BGP Sessions Start';
      sectionCls += '';
      break;
    default:
  }
  return (<p className={sectionCls}>{sectionTitle}</p>);
}


/**
 * RoutesLink is a link to the routes of the neighbor
 */
const RoutesLink = ({neighbor, children}) => {
  const { routeServerId } = useParams();
  if (!isUpState(neighbor.state)) {
    return <>{children}</>;
  };
  const url = `/routeservers/${routeServerId}/protocols/${neighbor.id}/routes`;
  return (
    <Link to={url}>{children}</Link>
  );
}

/**
 * Sort indicator indicates the sorting order of the column
 */
const SortIndicator = ({order, active}) => {
  if (!active) {
    return null;
  }
  let icon = faCircleArrowUp;
  if (order === 'desc') {
    icon = faCircleArrowDown;
  }
  return <FontAwesomeIcon icon={icon} />;
}


const ColumnHeader = ({title, column}) => {
  const [query, makeLink] = useQueryLink(querySortDefault);
  const sort = column.toLowerCase();
  const active = query.s === sort;

  let cls = `col-neighbor-attr col-neighbor-${column} `;
  let link = makeLink({s: sort}); // s: Sort column

  if (active) {
    cls += 'col-neighbor-active ';
    link = makeLink({o: query.o === 'asc' ? 'desc' : 'asc'}); // o: Toggle order
  }
  return (
    <th className={cls}>
      <Link to={link}>{title} <SortIndicator active={active} order={query.o} /></Link>
    </th>
  );
}

// Column Widgets:
const ColDescription = ({neighbor}) => {
  return (
    <td>
      <RoutesLink neighbor={neighbor}>
        {neighbor.description}
        {!isUpState(neighbor.state) && 
         neighbor.last_error &&
          <span className="protocol-state-error">
              {neighbor.last_error}
          </span>}
      </RoutesLink>
    </td>
  );
}

const ColUptime = ({neighbor}) => {
  return (
    <td className="date-since">
      <RelativeTimestamp value={neighbor.uptime} suffix={true} />
    </td>
  );
}


const ColLinked = ({neighbor, column}) => {
  // Access neighbor property by path
  const property = lookupProperty(neighbor, column);
  return (
    <td>
      <RoutesLink neighbor={neighbor}>
        {property}
      </RoutesLink>
    </td>
  );
}

const ColPlain = ({neighbor, column}) => {
  // Access neighbor property by path
  const property = useMemo(() => lookupProperty(neighbor, column), [neighbor, column]);
  return (
    <td>{property}</td>
  );
}

const ColNotAvailable = () => {
  return <td>-</td>;
}


const NeighborColumn = ({neighbor, column}) => {
  const rs = useRouteServer();
  const widgets = {
    // Special cases
    "asn": ColPlain,
    "state": ColPlain,

    "Uptime": ColUptime,
    "Description": ColDescription,
  };

  // For openbgpd the value is ommitted
  if (rs.type === "openbgpd") {
      widgets["routes_not_exported"] = ColNotAvailable;
  }

  // Get render function
  let Widget = widgets[column] || ColLinked;
  return (
    <Widget neighbor={neighbor} column={column} />
  );
}

const NeighborRow = ({neighbor, columns}) => {
  const cols = useMemo(() => columns.map((c) => 
    <NeighborColumn 
      key={c}
      neighbor={neighbor}
      column={c} />
  ), [neighbor, columns]);
  return (
    <tr>{cols}</tr>
  );
}

/**
 * NeighborsTable renders the table of neighbors
 */
const NeighborsTable = ({neighbors, state, ref}) => {
  const config = useConfig();
  const [query] = useQuery();

  const columns = config.neighbors_columns;
  const columnsOrder = config.neighbors_columns_order;

  const sortColumn = query.s;
  const sortOrder = query.o;

  const sortedNeighbors = useMemo(() =>
      sortNeighbors(neighbors, sortColumn, sortOrder
    ), [neighbors, sortColumn, sortOrder]);

  const header = useMemo(() => 
    columnsOrder.map((col) => 
      <ColumnHeader key={col} column={col} title={columns[col]} />
    ),
    [columns, columnsOrder]);


  if (!neighbors || neighbors.length === 0) {
    return null; // nothing to do here
  }

  const rows = sortedNeighbors.map((neighbor) => 
    <NeighborRow
      key={neighbor.id}
      columns={columnsOrder}
      neighbor={neighbor} />);

  return (
    <div className="card" ref={ref}>
      <Section state={state} />
      <table className="table table-striped table-protocols">
        <thead>
          <tr>
            {header}
          </tr>
        </thead>
        <tbody>
          {rows}
        </tbody>
      </table>
    </div>
  );
}

export default NeighborsTable;
