
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

import StartPage 
  from 'app/pages/start/Page';

import ConfigProvider
  from 'app/components/config/Provider';
import ContentProvider
  from 'app/components/content/Provider';

const Alice = () => {
  return (
    <ConfigProvider>
    <ContentProvider>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<StartPage />} />
      </Routes>
    </BrowserRouter>
    </ContentProvider>
    </ConfigProvider>
  );
}

export default Alice;

