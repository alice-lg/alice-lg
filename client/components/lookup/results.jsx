
import React from 'react'

export default class LookupResults extends React.Component {

    _countResults() {
        let count = 0;
        for (let rs in this.props.results) {
            let set = this.props.results[rs];
            count += set.length;
        }
        return count;
    }

    _resultSetEmpty() {
        let resultCount = this._countResults();
        if (this.props.finished && resultCount == 0){
            return true;
        }
        return false;
    }

    _awaitingResults() {
        let resultCount = this._countResults();
        if (!this.props.finished && resultCount == 0) {
            return true;
        }
        return false;
    }


    /* No Results */
    renderEmpty() {
        return (
            <div className="card card-results card-no-results">
                The prefix could not be found.
                Did you specify a network address?
            </div>
        );
    }

    render() {
        if (this._resultSetEmpty()) {
            return this.renderEmpty();
        }

        if (this._awaitingResults) {
            return null;
        }

        // Render Results table
        return (
            <div className="card card-results">
                ROUTES INCOMING!
            </div>
        );
    }

}


