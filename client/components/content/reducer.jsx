/*
 * Content reducer
 */

import {CONTENT_UPDATE} from './actions'

const initialState = {};

export default function reducer(state = initialState, action) {
  switch(action.type) {
    case CONTENT_UPDATE:
      return Object.assign({}, state, action.payload);
  }

  return state;
}

