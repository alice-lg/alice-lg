
import React from 'react'


export default class SearchInput extends React.Component {
  render() {
    return(
      <div className="input-group">
         <span className="input-group-addon">
          <i className="fa fa-search"></i>
         </span>
         <input type="text"
                className="form-control"
                {...this.props} />
      </div>
    );
  }
}


