
/**
 * DateTime component
 */

import { parseServerTime }
  from 'app/components/datetime/time';

/**
 * DateTime formats the provided datetime
 */
const DateTime = ({value, format="LLLL", utc=false}) => {
  let time = parseServerTime(value);
  if (utc) {
    time = time.utc();
  }
  return (<>{time.format(format)}</>);
}

export default DateTime;
