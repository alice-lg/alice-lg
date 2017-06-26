import React from 'react'
import DebounceInput from 'react-debounce-input'


export default class DebouncedSearchInput extends React.Component {
  render() {
    return(
      <div className="input-group">
         <span className="input-group-addon">
          <i className="fa fa-search"></i>
         </span>
         <DebounceInput
                minLength={2}
                debounceTimeout={250}
                className="form-control"
                {...this.props} />
      </div>
    );
  }
}



