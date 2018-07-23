
import React from 'react'
import {connect} from 'react-redux'

import LoadingIndicator
	from 'components/loading-indicator/small'


class Indicator extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      displayMessages: 0,
    };
  }

  isLoading(props) {
    return (props.receivedLoading ||
            props.filteredLoading ||
            props.notExportedLoading);
  }

  componentDidMount() {
    this.timeoutTimer = setInterval(
      () => this.tickMessages(), 1000
    );
  }

  componentWillUnmount() {
    clearInterval(this.timeoutTimer);
  }

  componentDidUpdate(prevProps) {
    if (!this.isLoading(this.props) &&
        this.isLoading(this.props) != this.isLoading(prevProps)) {
      // Stop timer
      this.setState({displayMessages: 0});
    }
  }

  tickMessages() {
    this.setState((prevState, props) => ({
      displayMessages: prevState.displayMessages + 1
    }));
  }

  render() {
    if (!this.isLoading(this.props)) {
      return null;
    }

    return (
      <div className="routes-loading card">
        <LoadingIndicator show={true} />

        {this.state.displayMessages >= 5 &&
          <p><br />&gt; Still loading routes, please be patient.</p>}
        {this.state.displayMessages >= 15 &&
          <p>&gt; This seems to take a while...</p>}
        {this.state.displayMessages >= 20 &&
          <p>&gt; This usually only happens when there are really many routes!<br />
             &nbsp; Please stand by a bit longer.</p>}

        {this.state.displayMessages >= 30 &&
          <p>&gt; This is taking really long...</p>}

        {this.state.displayMessages >= 40 &&
          <p>&gt; I heared there will be cake if you keep on waiting just a
             bit longer!</p>}

        {this.state.displayMessages >= 60 &&
          <p>&gt; I guess the cake was a lie.</p>}
      </div>
    );
  }

}


export default connect(
  (state) => ({
    receivedLoading:    state.routes.receivedLoading,
    filteredLoading:    state.routes.filteredLoading,
    notExportedLoading: state.routes.notExportedLoading,
  })
)(Indicator);


