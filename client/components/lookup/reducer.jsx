
import {SET_QUERY_TYPE,
        SET_QUERY_VALUE,
        SET_QUERY_INPUT_VALUE,

        LOOKUP_STARTED,
        LOOKUP_RESULTS,

        RESET,
        EXECUTE}
  from './actions'

import {QUERY_TYPE_UNKNOWN} from './query'

const initialState = {
	results: {},

    queue: new Set(),

	queryInput: "",

	query: "",
    queryType: QUERY_TYPE_UNKNOWN,

    queryRunning: false,
    queryFinished: false,
    queryDispatch: false,
};


// Action handlers:

// Handle lookup start
function _lookupStarted(state, lookup) {
    // Enqueue Routeserver
    let queue = new Set(state.queue);
    queue.add(lookup.routeserverId);

    // Clear results
    let results = Object.assign({}, state.results, {
        [lookup.routeserverId]: []
    });

    // Make state update
    return {
        queue: queue,
        results: results,

        queryRunning: true,
        queryFinished: false,
    };
}


// Handle a finished lookup
function _lookupResults(state, lookup) {
    // Dequeue routeserver
    let queue = new Set(state.queue);
    let currentQueueSize = queue.size;
    queue.delete(lookup.routeserverId);

    // Any routeservers left in the queue?
    let isRunning = true;
    if (queue.size == 0) {
        isRunning = false;
    }

    let isFinished = false;
    if (queue.size == 0 && currentQueueSize > 0) {
        isFinished = true;
    }


    // Update results set
    let results = Object.assign({}, state.results, {
        [lookup.routeserverId]: lookup.results,
    });

    // Make state update
    return {
        results: results,
        queue: queue,
        queryRunning: isRunning,
        queryFinished: isFinished
    }
}

// Reducer
export default function reducer(state=initialState, action) {
    let payload = action.payload;
	switch(action.type) {
        // Setup
        case SET_QUERY_TYPE:
		case SET_QUERY_VALUE:
        case SET_QUERY_INPUT_VALUE:
			return Object.assign({}, state, payload);

        // Search
        case LOOKUP_STARTED:
            // Update state on lookup started
            return Object.assign({}, state, _lookupStarted(state, payload), {
                queryDispatch: false,
            });

        case LOOKUP_RESULTS:
            // Update state when we receive results
            return Object.assign({}, state, _lookupResults(state, payload));

        case EXECUTE:
            return Object.assign({}, state, {
                queryDispatch: true,
            });

        case RESET:
            return Object.assign({}, state, initialState);
	}
	return state;
}

