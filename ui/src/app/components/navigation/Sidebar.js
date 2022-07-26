
import { Link }
  from 'react-router-dom';

import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faCloud }
  from '@fortawesome/free-solid-svg-icons';

import Content
  from 'app/components/content/Content';
import RouteServers
  from 'app/components/navigation/RouteServers';


const Sidebar = () => {
  return (
    <aside className="page-sidebar">
      <div className="sidebar-header">
        <div className="logo">
          <Link to='/'>
            <i>{/* Theme compatbility */}
            <FontAwesomeIcon 
              className="logo-icon"
              icon={faCloud} size="lg" transform="grow-11" />
            </i>
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
      <RouteServers />
    </aside>
  );  
}

export default Sidebar;
