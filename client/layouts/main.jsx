import React   from 'react'
import Sidebar from 'components/sidebar'

import ErrorsPage from 'components/errors/page'
import Config from 'components/config/view'

import Content from 'components/content'

export default class LayoutMain extends React.Component {
  render() {
    return (
      <div className="page">
        <ErrorsPage />
        <Sidebar />
        <div className="page-body">
          <main className="page-content">
            <div className="main-content-wrapper">
              {this.props.children}
            </div>
            <footer className="page-footer">
              <Content id="footer"></Content> 
            </footer>
          </main>
        </div>
        <Config/>
      </div>
    );
  }
}

