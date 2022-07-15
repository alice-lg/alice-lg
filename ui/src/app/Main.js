
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

import ErrorsProvider
  from 'app/components/errors/Provider';
import ConfigProvider
  from 'app/components/config/Provider';
import ContentProvider
  from 'app/components/content/Provider';
import RouteServersProvider
  from 'app/components/routeservers/Provider';

import Layout
  from 'app/components/page/Layout';
import RouteServer
  from 'app/components/routeservers/RouteServer';

import StartPage 
  from 'app/pages/StartPage';
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
            element={<RouteServer />}>

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
