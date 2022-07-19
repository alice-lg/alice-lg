
import { useState
       , useEffect
       }
  from 'react';

import LoadingIndicator
  from 'app/components/spinners/LoadingIndicator';


const WaitingText = ({resource = "routes"}) => {
  const [time, setTime] = useState(0);

  useEffect(() => {
    const tRef = setInterval(() => {
      setTime((t) => (t += 1));
    }, 1000);
    return () => {
      clearInterval(tRef);
    }
  }, []);

  return (
    <div className="routes-loading card">
      <LoadingIndicator show={true} />

      {time >= 5 &&
        <p><br />&gt; Still loading routes, please be patient.</p>}
      {time >= 15 &&
        <p>&gt; This seems to take a while...</p>}
      {time >= 20 &&
        <p>&gt; This usually only happens when there are really many routes!<br />
           &nbsp; Please stand by a bit longer.</p>}

      {time >= 30 &&
        <p>&gt; This is taking really long...</p>}

      {time >= 40 &&
        <p>&gt; I heard there will be cake if you keep on waiting just a
           bit longer!</p>}

      {time >= 60 &&
        <p>&gt; I guess the cake was a lie.</p>}
    </div>
  );
}

export default WaitingText;
