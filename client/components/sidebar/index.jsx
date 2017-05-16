
/**
 * Main Sidebar Component (aka. Navigation)
 *
 */


import React from 'react'

import Header from './header'
import Routeservers from './routeservers'

export default class Sidebar extends React.Component {

  render() {
    return (
      <aside className="page-sidebar">
        <Header />
        <Routeservers />
      </aside>
    )
  }

}


