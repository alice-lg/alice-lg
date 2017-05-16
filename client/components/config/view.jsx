import React from 'react'
import { connect } from 'react-redux'

import { loadConfig } from 'components/config/actions'


class Config extends React.Component {
  componentDidMount() {
    this.props.dispatch(loadConfig());
  }

  render() {
    return null;
  }
}

export default connect()(Config);
