
import { useState
       , useEffect
       }
  from 'react';


const WaitingText = ({resource}) => {
  const [time, setTime] = useState(0);

  useEffect(() => {
    const tRef = setInterval(() => {
      setTime((t) => (t += 1));
    }, 1000);
    return () => {
      clearInterval(tRef);
    }
  }, []);

  if (time < 5) {
    return null;
  }

  return (
    <div className="card routes-loading">
      {time >= 5 &&
        <p>&gt; Still loading routes, please be patient.</p>}
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

const WaitingCard = ({
  isLoading,
  resource = "routes",
}) => {
  if (!isLoading) {
    return null;
  }
  return <WaitingText />
}

export default WaitingCard;
