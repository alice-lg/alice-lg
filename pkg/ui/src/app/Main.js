
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

import Page
  from 'app/components/page/Page';

import StartPage 
  from 'app/pages/start/Page';
import NeighborsPage
  from 'app/pages/neighbors/Page';
import NotFoundPage
  from 'app/pages/errors/NotFound';



const Main = () => {
  return (
    <ErrorsProvider>
    <ConfigProvider>
    <RouteserversProvider>
    <ContentProvider>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Page />}>
          <Route index element={<StartPage />} />
          <Route path="routeservers/:routeServerId">
            <Route index element={<NeighborsPage />} />
          </Route>
          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
    </ContentProvider>
    </RouteserversProvider>
    </ConfigProvider>
    </ErrorsProvider>
  );
}

export default Main;

