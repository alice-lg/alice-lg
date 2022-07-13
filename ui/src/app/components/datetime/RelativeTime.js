
import moment from 'moment';

import { useEffect
       , useState
       }
  from 'react';


/**
 * RelativeTime renders the relative time
 */
const RelativeTime = ({
  value, 
  suffix,
  fuzzyNow=false,
  pastEvent=false,
  futureEvent=false
}) => {
  const [now, setNow] = useState(moment.utc());

  // Update current time
  useEffect(() => {
    const tRef = setInterval(() => setNow(moment.utc()), 1000);
    return () => {
      clearInterval(tRef);
    };
  }, [])

  // Time can be capped, if we are handling a past
  // or future event.
  const capTime = (t) => {
    if (pastEvent && t.isAfter(now)) {
      return now;
    }
    if (futureEvent && t.isBefore(now)) {
      return now;
    }
    return t;
  }

  if (!value) {
    return null;
  }

  let time = value;
  if (!(value instanceof moment)) {
    time = moment.utc(value); 
  }
  time = capTime(time);

  if (fuzzyNow) {
    if (Math.abs(now - time) / 1000.0 < fuzzyNow) {
      return (
        <>just now</>
      );
    }
  }

  return (
    <>{time.fromNow(suffix)}</>
  );
}

export default RelativeTime;
