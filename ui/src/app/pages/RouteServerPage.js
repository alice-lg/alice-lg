
import { Outlet }
  from 'react-router-dom';

import { useParams }
  from 'react-router-dom';

import { RouteServerStatusProvider }
  from 'app/context/route-servers';
import { NeighborsProvider }
  from 'app/context/neighbors';

/**
 * The RouteServerPage component initializes the routeserver status
 * context and the neighbors context.
 */
const RouteServerPage = () => {
  const { routeServerId } = useParams();
  return (
    <RouteServerStatusProvider routeServerId={routeServerId}>
    <NeighborsProvider routeServerId={routeServerId}>
      <Outlet />
    </NeighborsProvider>
    </RouteServerStatusProvider>
  );
}

export default RouteServerPage;
