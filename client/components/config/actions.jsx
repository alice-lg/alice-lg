import axios from 'axios';
import {apiError} from 'components/errors/actions'
import {loadRejectReasonsSuccess} from 'components/routeservers/actions';

export const LOAD_CONFIG_SUCCESS = "@birdseye/LOAD_CONFIG_SUCCESS";

function loadConfigSuccess(routes_columns) {
  return {
    type: LOAD_CONFIG_SUCCESS,
    routes_columns: routes_columns
  }
}

export function loadConfig() {
  return (dispatch) => {
    axios.get(`/api/config`)
      .then(({data}) => {
        dispatch(
          loadRejectReasonsSuccess(data.config.rejection.asn,
                                   data.config.rejection.reject_id,
                                   data.config.reject_reasons)
        );
        dispatch(loadConfigSuccess(data.config.routes_columns));
      })
      .catch(error => dispatch(apiError(error)));
  }
}
