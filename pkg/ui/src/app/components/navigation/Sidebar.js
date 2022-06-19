
import { Link }
  from 'react-router-dom';

import Content
  from 'app/components/content/Content';

const Sidebar = () => {
  return (
    <aside className="page-sidebar">
      <div className="sidebar-header">
        <div className="logo">
          <Link to='/'>
            <i className="fa fa-cloud"></i>
          </Link>
        </div>
        <div className="title">
          <h1><Content id="header.title">Alice</Content></h1>
          <p>
            <Content id="header.tagline">
              Your friendly BGP looking glass.
            </Content>
          </p>
        </div>
      </div>


      RouteServers
    </aside>
  );  
}

export default Sidebar;
