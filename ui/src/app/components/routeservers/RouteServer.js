
import { Outlet }
  from 'react-router-dom';

import { useParams }
  from 'react-router-dom';

import RouteServerStatusProvider
  from 'app/components/routeservers/StatusProvider';
import NeighborsProvider
  from 'app/components/neighbors/Provider';

/**
 * The RouteServer component initializes the routeserver status
 * context and the neighbors context.
 */
const RouteServer = () => {
  const { routeServerId } = useParams();
  return (
    <RouteServerStatusProvider routeServerId={routeServerId}>
    <NeighborsProvider routeServerId={routeServerId}>
      <Outlet />
    </NeighborsProvider>
    </RouteServerStatusProvider>
  );
}

export default RouteServer;
