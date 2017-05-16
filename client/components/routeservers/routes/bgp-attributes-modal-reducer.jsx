
/**
 * Bgp Attributes Modal Reducer
 *
 * @author Matthias Hannig <mha@ecix.net>
 */


import {SHOW_BGP_ATTRIBUTES_MODAL,
        HIDE_BGP_ATTRIBUTES_MODAL,
        SET_BGP_ATTRIBUTES} from './bgp-attributes-modal-actions'

const initialState = {
  show: false,
  bgpAttributes: {}
};


export default function reducer(state = initialState, action) {
  switch(action.type) {
    case SHOW_BGP_ATTRIBUTES_MODAL:
      return Object.assign({}, state, { show: true });
    case HIDE_BGP_ATTRIBUTES_MODAL:
      return Object.assign({}, state, { show: false });
    case SET_BGP_ATTRIBUTES:
      return Object.assign({}, state, action.payload);
  }

  return state;
}

