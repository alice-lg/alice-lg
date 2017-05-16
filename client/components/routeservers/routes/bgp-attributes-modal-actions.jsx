

export const SHOW_BGP_ATTRIBUTES_MODAL = '@birdseye/SHOW_BGP_ATTRIBUTES_MODAL';
export const HIDE_BGP_ATTRIBUTES_MODAL = '@birdseye/HIDE_BGP_ATTRIBUTES_MODAL';
export const SET_BGP_ATTRIBUTES = '@birdseye/SET_BGP_ATTRIBUTES';


/**
 * Action Creators
 */

export function showBgpAttributesModal() {
  return {
    type: SHOW_BGP_ATTRIBUTES_MODAL,
  }
}

export function hideBgpAttributesModal() {
  return {
    type: HIDE_BGP_ATTRIBUTES_MODAL
  }
}


export function setBgpAttributes(attributes) {
  return {
    type: SET_BGP_ATTRIBUTES,
    payload: {
      bgpAttributes: attributes
    }
  }
}

export function showBgpAttributes(attributes) {
  return (dispatch) => {
    dispatch(setBgpAttributes(attributes));
    dispatch(showBgpAttributesModal());
  }
}

