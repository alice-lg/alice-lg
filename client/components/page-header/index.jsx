
import React from 'react'


export default class PageHeader extends React.Component {
  render() {
    return (
      <div className="page-header">
        {this.props.children}
      </div>
    )
  }
}

