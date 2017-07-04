import axios from 'axios';
import {apiError} from 'components/errors/actions'
import {loadRejectReasonsSuccess,
        loadNoExportReasonsSuccess}
  from 'components/routeservers/large-communities/actions';

export const LOAD_CONFIG_SUCCESS = "@birdseye/LOAD_CONFIG_SUCCESS";

function loadConfigSuccess(config) {
  return {
    type: LOAD_CONFIG_SUCCESS,
    payload: config
  };
}

export function loadConfig() {
  return (dispatch) => {
    axios.get(`/api/config`)
      .then(({data}) => {
        dispatch(
          loadRejectReasonsSuccess(data.rejection.asn,
                                   data.rejection.reject_id,
                                   data.reject_reasons)
        );
        dispatch(
            loadNoExportReasonsSuccess(
                data.noexport.asn,
                data.noexport.noexport_id,
                data.noexport_reasons)
        );
        dispatch(loadConfigSuccess(data));
      })
      .catch(error => dispatch(apiError(error)));
  }
}
