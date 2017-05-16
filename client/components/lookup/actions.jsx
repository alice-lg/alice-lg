
import axios from 'axios'

export const SET_QUERY_INPUT_VALUE = "@lookup/SET_QUERY_INPUT_VALUE";
export const SET_QUERY_VALUE       = "@lookup/SET_QUERY_VALUE";
export const SET_QUERY_TYPE        = "@lookup/SET_QUERY_TYPE";

export const RESET = "@lookup/RESET";
export const EXECUTE = "@lookup/EXECUTE";

export const LOOKUP_STARTED = "@lookup/LOOKUP_STARTED";
export const LOOKUP_RESULTS = "@lookup/LOOKUP_RESULTS";


/*
 * Action Creators
 */

export function setQueryInputValue(q) {
    if(!q) { q = ''; }
	return {
		type: SET_QUERY_INPUT_VALUE,
		payload: {
			queryInput: q
		}
	}
}

export function setQueryValue(q) {
	return {
		type: SET_QUERY_VALUE,
		payload: {
			query: q
		}
	}
}

export function setQueryType(type) {
    return {
        type: SET_QUERY_TYPE,
        payload: {
            queryType: type
        }
    }
}

export function reset() {
    return {
        type: RESET
    }
}

export function execute() {
    return {
        type: EXECUTE
    }
}


export function lookupStarted(routeserverId, query) {
    return {
        type: LOOKUP_STARTED,
        payload: {
            routeserverId: routeserverId,
            query: query
        }
    }
}


export function lookupResults(routeserverId, query, results) {
    return {
        type: LOOKUP_RESULTS,
        payload: {
            routeserverId: routeserverId,
            query: query,
            results: results
        }
    }
}



export function routesSearch(routeserverId, q) {
    return (dispatch) => {
        dispatch(lookupStarted(routeserverId, q));
        axios.get(`/birdseye/api/routeserver/${routeserverId}/routes/lookup?q=${q}`)
            .then((result) => {
                let routes = result.data.result.routes;
                dispatch(lookupResults(
                    routeserverId,
                    q,
                    routes
                ));
            })
            .catch((error) => {
                dispatch(lookupResults(
                    routeserverId,
                    q,
                    []
                ));
            });
    }
}

