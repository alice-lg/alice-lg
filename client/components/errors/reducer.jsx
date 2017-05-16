import {API_ERROR} from './actions'

const initialState = {
  error: null,
};


export default function reducer(state = initialState, action) {
  switch(action.type) {
    case API_ERROR:
      return {error: action.error};
  }
  return state;
}



