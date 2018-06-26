
import React from 'react'
import {Link} from 'react-router'

import content from 'helpers/content'

export default class SidebarHeader extends React.Component {
  render() {
    return (
      <div className="sidebar-header">
        <div className="logo">
          <Link to='/'>
            <i className={content("header.icon", "fa fa-cloud")}></i>
          </Link>
        </div>
        <div className="title">
          <h1>{content("header.title", "Alice")}</h1>
          <p>{content("header.tagline",
                      "Your friendly bird looking glass")}</p>
        </div>
      </div>
    );
  }
}

