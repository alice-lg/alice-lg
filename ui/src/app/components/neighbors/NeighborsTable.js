
import bigInt from 'big-integer';
import { ipToNumeric } from 'lib/ip';

import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faCircleArrowUp
       , faCircleArrowDown
       , faCirclePlus
       , faCircleMinus
       }
  from '@fortawesome/free-solid-svg-icons';


import { useMemo, useState }
  from 'react';

import { useParams
       , Link
       }
  from 'react-router-dom';


import { useConfig }
  from 'app/context/config';
import { useRouteServer }
  from 'app/context/route-servers';
import { useQuery
       , useQueryLocation
       , PARAM_SORT
       , PARAM_ORDER
       }
  from 'app/context/query';

import { isUpState } 
  from 'app/components/neighbors/state';
import RelativeTimestamp
  from 'app/components/datetime/RelativeTimestamp';
import AsnLink
  from 'app/components/asns/AsnLink';

/**
 * Default: Sort by ASN, ascending order.
 */
const querySortDefault = {
  [PARAM_SORT]: 'asn',
  [PARAM_ORDER]: 'asc',
};

const lookupProperty = (obj, path) => {
  let property = path.split(".").reduce((acc, part) => acc[part], obj);
  if (typeof(property) == "undefined") {
    property = `Property "${path}" not found in object.`;
  }

  return property;
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
    case 'idle':
      sectionTitle = 'Idle BGP Sessions';
      sectionCls += 'idle ';
      break;
    case 'down':
      sectionTitle = 'BGP Sessions Down';
      sectionCls += 'down ';
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
  const url = `/routeservers/${routeServerId}/neighbors/${neighbor.id}/routes`;
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


/**
 * Column Header with sorting indicator
 */
const ColumnHeader = ({title, column}) => {
  const sort = column.toLowerCase();
  const [query] = useQuery(querySortDefault);
  const columnSort = useQueryLocation({
    ...query,
    [PARAM_SORT]: sort,
  });
  const toggleOrder = useQueryLocation({
    ...query,
    [PARAM_ORDER]: query[PARAM_ORDER] === 'asc' ? 'desc' : 'asc',
  });
  
  const active = query[PARAM_SORT] === sort;
  let cls = `col-neighbor-attr col-neighbor-${column} `;
  let link = columnSort; 
  if (active) {
    cls += 'col-neighbor-active ';
    link = toggleOrder;
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

const ColAsn = ({neighbor}) => {
  return (
    <td><AsnLink asn={neighbor.asn} /></td>
  );
}

const ColNotAvailable = () => {
  return <td>-</td>;
}


const ColToggleExtended = ({extended, neighbor}) => {
    const [isExtended, setExtended] = extended;
    let icon = faCirclePlus;
    if (isExtended) {
        icon = faCircleMinus;
    }
    const channels = neighbor?.routes_channels;
    const isMultiChannel = (channels?.ipv4 && channels?.ipv6);
    if (!isMultiChannel) {
        return <td></td>;
    }

    return (
      <td>
        <button
          className="btn btn-xs btn-extend"
          onClick={() => setExtended(!isExtended)}>
            <FontAwesomeIcon icon={icon} />
        </button>
      </td>
    );
}


const NeighborColumn = ({neighbor, column}) => {
  const rs = useRouteServer();
  const widgets = {
    // Special cases
    "asn": ColAsn,
    "state": ColPlain,

    "Uptime": ColUptime,
    "Description": ColDescription,
  };

  // For openbgpd the value is omitted
  if (rs.type === "openbgpd") {
      widgets["routes_not_exported"] = ColNotAvailable;
  }

  // Get render function
  let Widget = widgets[column] || ColLinked;
  return (
    <Widget neighbor={neighbor} column={column} />
  );
}


const NeighborRow = ({neighbor, columns, extended}) => {
  const cols = useMemo(() => columns.map((c) =>
    <NeighborColumn
      key={c}
      neighbor={neighbor}
      column={c} />
  ), [neighbor, columns]);
  return (
    <tr>
        {cols}
        <ColToggleExtended neighbor={neighbor} extended={extended} />
    </tr>
  );
}

const NeighborRowDetails = ({neighbor, columns, channel, extended}) => {
    const [isExtended] = extended;
    if (!isExtended || !channel) {
        return null;
    }

    let chan = neighbor?.routes_channels?.ipv4;
    if (channel === 6) {
        chan = neighbor?.routes_channels?.ipv6;
    }

    const cols = columns.map((c) => {
        let cval = chan[c];
        if (cval) {
            return <td key={c}>{cval}</td>;
        }
        if (c === "Description") {
            return <td key={c} className="col-ip-info">IPv{channel}</td>;
        }
        return <td key={c}></td>;
    });

    return <tr>{cols}<td></td></tr>;
}

const TableNeighbor = ({columns, neighbor}) => {
    const extended = useState(false);
    return <>
        <NeighborRow
            columns={columns}
            neighbor={neighbor}
            extended={extended} />
        <NeighborRowDetails
            columns={columns}
            neighbor={neighbor}
            channel={4}
            extended={extended} />
        <NeighborRowDetails
            columns={columns}
            neighbor={neighbor}
            channel={6}
            extended={extended} />
    </>;
};

/**
 * NeighborsTable renders the table of neighbors
 */
const NeighborsTable = ({neighbors, state, ref}) => {
  const config = useConfig();
  const [query] = useQuery();

  const columns = config.neighbors_columns;
  const columnsOrder = config.neighbors_columns_order;

  const sortColumn = query[PARAM_SORT];
  const sortOrder = query[PARAM_ORDER];

  const sortedNeighbors = useMemo(() =>
      sortNeighbors(
        neighbors,
        sortColumn,
        sortOrder,
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
    <TableNeighbor
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
