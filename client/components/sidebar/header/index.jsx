
import React from 'react'
import {Link} from 'react-router'

export default class SidebarHeader extends React.Component {
  render() {
    return (
      <div className="sidebar-header">
        <div className="logo">
          <Link to='/'>
            <i className="fa fa-cloud"></i>
          </Link>
        </div>
        <div className="title">
          <h1>Birdseye</h1>
          <p>Your friendly bird looking glass</p>
        </div>
      </div>
    );
  }
}

