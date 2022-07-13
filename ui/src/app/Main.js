
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
import RouteserversProvider
  from 'app/components/routeservers/Provider';

import Layout
  from 'app/components/page/Layout';

import StartPage 
  from 'app/pages/StartPage';
import NeighborsPage
  from 'app/pages/NeighborsPage';
import NotFoundPage
  from 'app/pages/NotFoundPage';

const Main = () => {
  return (
    <ErrorsProvider>
    <ConfigProvider>
    <RouteserversProvider>
    <ContentProvider>
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route index element={<StartPage />} />

          {/* RouteServers */}
          <Route path="routeservers/:routeServerId">
            <Route index element={<NeighborsPage />} />
          </Route>

          {/* Fallback */}
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </Layout>
    </BrowserRouter>
    </ContentProvider>
    </RouteserversProvider>
    </ConfigProvider>
    </ErrorsProvider>
  );
}

export default Main;

