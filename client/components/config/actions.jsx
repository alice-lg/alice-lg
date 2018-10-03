import axios from 'axios';
import {apiError} from 'components/errors/actions'

export const LOAD_CONFIG_SUCCESS = "@config/LOAD_CONFIG_SUCCESS";

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
        dispatch(loadConfigSuccess(data));
      })
      .catch(error => dispatch(apiError(error)));
  }
}
