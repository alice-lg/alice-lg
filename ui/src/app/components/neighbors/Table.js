
/**
 * NeighborsTable renders the table of neighbors
 */
const NeighborsTable = ({neighbors, state, ref}) => {

  if (!neighbors || neighbors.length === 0) {
    return null; // nothing to do here
  }

  let sectionTitle = '';
  let sectionAnchor = 'sessions-unknown';
  let sectionCls = 'card-header card-header-neighbors ';

  switch(state) {
    case 'up':
      sectionAnchor = 'sessions-up';
      sectionTitle  = 'BGP Sessions Established';
      sectionCls  += 'established ';
      break;
    case 'down':
      sectionAnchor = 'sessions-down';
      break;
    case 'start':
      sectionAnchor = 'sessions-down';
      sectionTitle = 'BGP Sessions Down';
      sectionCls += 'down ';
      break;
    default:
  }

  let header = <td>Header</td>;
  let rows = <tr><td>Row</td></tr>;

  return (
    <div className="card" ref={ref}>
      <p className={sectionCls}>{sectionTitle}</p>

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
