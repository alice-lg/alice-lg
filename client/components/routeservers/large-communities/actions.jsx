
export const LOAD_REJECT_REASONS_REQUEST = '@birdseye/LOAD_REJECT_REASONS_REQUEST';
export const LOAD_REJECT_REASONS_SUCCESS = '@birdseye/LOAD_REJECT_REASONS_SUCCESS';

export const LOAD_NOEXPORT_REASONS_REQUEST = '@birdseye/LOAD_NOEXPORT_REASONS_REQUEST';
export const LOAD_NOEXPORT_REASONS_SUCCESS = '@birdseye/LOAD_NOEXPORT_REASONS_SUCCESS';


export function loadRejectReasonsSuccess(asn, reject_id, reject_reasons) {
  return {
    type: LOAD_REJECT_REASONS_SUCCESS,
    payload: {
        reject_asn: asn,
        reject_id: reject_id,
        reject_reasons: reject_reasons}
  };
}


export function loadNoExportReasonsSuccess(asn, noexport_id, reasons) {
  return {
    type: LOAD_NOEXPORT_REASONS_SUCCESS,
    payload: {
       noexport_asn: asn,
       noexport_id: noexport_id,
       noexport_reasons: reasons
    }
  };
}


