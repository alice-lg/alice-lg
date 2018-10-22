

// Actions
export const SET_FILTER_VALUE = "@neighbors/SET_FILTER_VALUE";

// Action Creators: Set Filter Query
export function setFilterValue(value) {
  return {
    type: SET_FILTER_VALUE,
    payload: {
      value: value
    }
  }
}

