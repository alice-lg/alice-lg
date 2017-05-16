
/**
 * Bootstrap Modal React Component
 *
 * @author Matthias Hannig <mha@ecix.net>
 */

import React from 'react'

export class Header extends React.Component {
  render() {
    return(
      <div className="modal-header">
        <button type="button"
                className="close" 
                aria-label="Close"
                onClick={this.props.onClickClose}>
                 <span aria-hidden="true">&times;</span></button>

        {this.props.children}
      </div>
    );
  }
}

export class Body extends React.Component {
  render() {
    return (
      <div className="modal-body">
        {this.props.children}
      </div>
    );
  }
}


export class Footer extends React.Component {
  render() {
    return(
      <div className="modal-footer">
        {this.props.children}
      </div>
    );
  }
}

export default class Modal extends React.Component {
  render() {
    if(!this.props.show) {
      return null;
    }

    return (
      <div className={this.props.className}>
        <div className="modal modal-open modal-show fade in" role="dialog">
          <div className="modal-dialog" role="document">
            <div className="modal-content">
              {this.props.children}
            </div>
          </div>
        </div>
        <div className="modal-backdrop fade in"
             onClick={this.props.onClickBackdrop}></div>
      </div>
    );

  }
}


