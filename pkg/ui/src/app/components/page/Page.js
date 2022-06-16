/*
 * Page implements the default alice page with
 * the routeserver navigation on the left, a header on top
 * and a content view in the middle.
 */

import Content
  from 'app/components/content/Content';

// TODO
const Errors = () => <></>;

const Page = ({children}) => {

  // Main Layout
  return (
    <div className="page">
      <Errors />
      <aside className="page-sidebar">
        SidebarHeader
        RouteServers
      </aside>
      <div className="page-body">
        <main className="page-content">
          <div className="main-content-wrapper">
            {children}
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
