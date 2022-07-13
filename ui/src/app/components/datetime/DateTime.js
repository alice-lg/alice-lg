
/**
 * DateTime component
 */

import { parseServerTime }
  from 'app/components/datetime/time';

/**
 * DateTime formats the provided datetime
 */
const DateTime = ({value, format="LLLL"}) => {
  const time = parseServerTime(value);
  return (<>{time.format(format)}</>);
}

export default DateTime;
