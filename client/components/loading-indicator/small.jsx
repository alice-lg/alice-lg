
import React from 'react'
import Spinner from 'react-spinkit'

export default class Indicator extends React.Component {
	render() {
		if (this.props.show == false) {
			return null;
		}

		return (
			<div className="loading-indicator">
       	<Spinner spinnerName="circle" />
			</div>
		);
	}
}


