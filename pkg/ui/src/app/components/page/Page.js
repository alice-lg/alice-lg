/*
 * Page implements the default alice page with
 * the routeserver navigation on the left, a header on top
 * and a content view in the middle.
 */

import { Outlet }
  from 'react-router-dom';

import Content
  from 'app/components/content/Content';
import Errors
  from 'app/components/errors/Errors';
import NavigationSidebar
  from 'app/components/navigation/Sidebar';


const Page = ({children}) => {

  // Main Layout
  return (
    <div className="page">
      <Errors />
      <NavigationSidebar />
      <div className="page-body">
        <main className="page-content">
          <div className="main-content-wrapper">
            <Outlet /> 
          </div>
          <footer className="page-footer">
            <Content id="footer"></Content> 
          </footer>
        </main>
      </div>
    </div>
  );
}

export default Page;
