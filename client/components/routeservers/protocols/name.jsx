
import React from 'react'
import {connect} from 'react-redux'

/*
 * Show current neighbour if selected
 * This should help to create a breadcrumb style navigation
 * in the header.
 */
class ProtocolName extends React.Component {
    render() {
        return (
            <span className="status-protocol">
                {this.props.protocol.description}
            </span>
        );
    }
}

export default connect(
    (state, props) => {
        let rsProtocols = state.routeservers.protocols[props.routeserverId]||[];
        let protocol = rsProtocols.filter((p) => {
            return p.id == props.protocolId;
        })[0]||{};
        return {
            protocol: protocol
        };
    }
)(ProtocolName);
