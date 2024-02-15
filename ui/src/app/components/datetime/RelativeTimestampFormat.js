
import moment from 'moment'

/**
 * Render the formatted 'absolute' time when given a
 * relative timestamp (in nanoseconds).
 *
 * The timestamp is the duration from now to the absolute
 * date time in the past.
 */
const RelativeTimestampFormat = ({value, format, now}) => {
  if (!now) {
    now = moment().utc();
  } else {
    now = moment(now);
  }
  const tsMs = value / 1000.0 / 1000.0; // nano -> micro -> milli
  const abs = now.subtract(tsMs, 'ms');
  return (
    <>{abs.format(format)}</>
  );
}

export default RelativeTimestampFormat;
