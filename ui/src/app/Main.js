
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
import SearchGlobalPage 
  from 'app/pages/SearchGlobalPage';
import NotFoundPage
  from 'app/pages/NotFoundPage';


const Routing = () => (
  <Routes>
    <Route index element={<StartPage />} />

    {/* RouteServers */}
    <Route
      path="routeservers/:routeServerId"
      element={<RouteServerPage />}>

      <Route index element={<NeighborsPage />} />

      {/* Neighbors */}
      <Route
        path="neighbors/:neighborId/routes"
        element={<RoutesPage />} />
      {/* DEPRECATION NOTICE: The 'protocols' route will be */}
      {/*   removed and is only here for backwards compatibility */}
      <Route
        path="protocols/:neighborId/routes"
        element={<RoutesPage />} />

    </Route>
  
    {/* Search */}
    <Route path="search" element={<SearchGlobalPage />} />

    {/* Fallback */}
    <Route path="*" element={<NotFoundPage />} />
  </Routes>
);

const Main = () => {
  return (
    <ErrorsProvider>
    <ConfigProvider>
    <RouteServersProvider>
    <ContentProvider>
    <BrowserRouter>
      <Layout>
        <Routing />
      </Layout>
    </BrowserRouter>
    </ContentProvider>
    </RouteServersProvider>
    </ConfigProvider>
    </ErrorsProvider>
  );
}

export default Main;
