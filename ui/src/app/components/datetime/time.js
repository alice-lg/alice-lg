
import moment from 'moment';

const SERVER_TIME_FMT = "YYYY-MM-DDTHH:mm:ss.SSSSSSSSZ";

export const parseServerTime = (serverTime) => {
  return moment(serverTime, SERVER_TIME_FMT);
}

