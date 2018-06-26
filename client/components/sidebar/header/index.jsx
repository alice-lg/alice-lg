
import React from 'react'
import {Link} from 'react-router'

import Content from 'components/content'


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
          <h1><Content id="header.title">Alice</Content></h1>
          <p>
            <Content id="header.tagline">
              Your friendly bird looking glass.
            </Content>
          </p>
        </div>
      </div>
    );
  }
}

