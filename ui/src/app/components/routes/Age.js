
import { useMemo }
  from 'react';

import { useApiStatus }
  from 'app/context/api-status';


import RelativeTimestampFormat
  from 'app/components/datetime/RelativeTimestampFormat';
import RelativeTimestamp
  from 'app/components/datetime/RelativeTimestamp';

export const RouteAgeDetails = ({route}) => {
  const api = useApiStatus();

  return useMemo(() =>
    <>
      <RelativeTimestampFormat
        value={route.age}
        now={api.receivedAt}
        format="YYYY-MM-DD HH:mm:ss"/> UTC
        <b> (<RelativeTimestamp
            value={route.age}
            now={api.receivedAt}
            suffix={true} />)
        </b>
    </>,
    [route.age, api.receivedAt]
  );
}

export const RouteAgeRelative = ({route}) => {
  const api = useApiStatus();

  return useMemo(() =>
    <RelativeTimestamp
      value={route.age}
      now={api.receivedAt}
      suffix={true} />,
    [route.age, api.receivedAt]
  );
}

export const RouteAgeAbsolute = ({route}) => {
  const api = useApiStatus();

  return useMemo(() =>
    <><RelativeTimestampFormat
      value={route.age}
      now={api.receivedAt}
      format="YYYY-MM-DD HH:mm:ss"/> UTC</>
    , [route.age, api.receivedAt]
  );
}

