import moment from 'moment'

/**
 * Render a relative timestamp
 */
const RelativeTimestamp = ({value, suffix}) => {
  const tsMs = value / 1000.0 / 1000.0; // nano -> micro -> milli
  const now = moment.utc()
  const rel = now.subtract(tsMs, 'ms');
  return (
    <>{rel.fromNow(suffix)}</>
  );
}

export default RelativeTimestamp;
