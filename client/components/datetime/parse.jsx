
/*
 * Some datetime parsing helper functions
 */

import moment from 'moment'
 

window.moment = moment;

export function parseServerTime(serverTime) {
  const fmt = "YYYY-MM-DDTHH:mm:ss.SSSSSSSSZ"; // S was 4 byte short
  return moment(serverTime, fmt);
}


