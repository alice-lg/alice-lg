
import React from 'react'
import {connect} from 'react-redux'

import SearchInput
  from 'components/search-input'

import LoadingIndicator
	from 'components/loading-indicator/small'

import {setQueryInputValue,
        execute,
        routesSearch}
	from './actions'


import QueryDispatcher
  from './query-dispatcher'

import LookupResults
  from './results'

import {queryParams}
  from 'components/utils/query'


class LookupView extends React.Component {

	setQuery(q) {
		this.props.dispatch(
			setQueryInputValue(q)
		);
	}

    componentDidMount() {
        // Initial mount: keep query from querystring
        let params = queryParams();
		this.props.dispatch(
			setQueryInputValue(params.q)
		);
    }


    handleFormSubmit(e) {
        e.preventDefault();
        this.props.dispatch(execute());
        return false;
    }

	render() {
		return (
			<div className="routes-lookup">

				<div className="card lookup-header">
                    <form className="form-lookup" onSubmit={(e) => this.handleFormSubmit(e)}>
                        <SearchInput placeholder="Search for routes by entering a network address"
                                     name="q"
                                     onChange={(e) => this.setQuery(e.target.value)}
                                     disabled={this.props.isSearching}
                                     value={this.props.queryInput} />
                        <QueryDispatcher />
                    </form>
				</div>
				<LoadingIndicator show={this.props.isRunning} />
				<div className="lookup-results">
                    <LookupResults results={this.props.results}
                                   finished={this.props.isFinished} />
				</div>
			</div>
		);
	}
}

export default connect(
	(state) => {
		return {
			isRunning: state.lookup.queryRunning,
            isFinished: state.lookup.queryFinished,

            queryInput: state.lookup.queryInput,

			results: state.lookup.results,
            search: state.lookup.search,
		}
	}
)(LookupView);

