
/**
 * Alice (formerly known as Birdseye) UI
 * -------------------------------------
 *
 * @author Annika Hannig <annika@hannig.cc>
 */

import { BrowserRouter
       , Routes
       , Route
       } 
  from 'react-router-dom';

import { ErrorsProvider }
  from 'app/context/errors';
import { ConfigProvider }
  from 'app/context/config';
import { ContentProvider }
  from 'app/context/content';
import { RouteServersProvider }
  from 'app/context/route-servers';

import Layout
  from 'app/components/page/Layout';

import StartPage 
  from 'app/pages/StartPage';
import RouteServerPage
  from 'app/pages/RouteServerPage';
import NeighborsPage
  from 'app/pages/NeighborsPage';
import RoutesPage
  from 'app/pages/RoutesPage';
import NotFoundPage
  from 'app/pages/NotFoundPage';

const Main = () => {
  return (
    <ErrorsProvider>
    <ConfigProvider>
    <RouteServersProvider>
    <ContentProvider>
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route index element={<StartPage />} />

          {/* RouteServers */}
          <Route path="routeservers/:routeServerId"
            element={<RouteServerPage />}>

            <Route index
              element={<NeighborsPage />} />
            <Route path="protocols/:neighborId/routes"
              element={<RoutesPage />} />

          </Route>

          {/* Fallback */}
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </Layout>
    </BrowserRouter>
    </ContentProvider>
    </RouteServersProvider>
    </ConfigProvider>
    </ErrorsProvider>
  );
}

export default Main;
