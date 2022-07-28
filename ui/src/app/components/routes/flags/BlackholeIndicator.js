
import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faCircle }
  from '@fortawesome/free-solid-svg-icons';

import { useRouteServer }
  from 'app/context/route-servers';


const BlackholeIndicator = ({route}) => {
  const routeServer = useRouteServer(); // blackholes are store per RS

  const blackholes = routeServer?.blackholes || [];
  const communities = route?.bgp?.communities || [];
  const nextHop = route?.bgp?.next_hop;

  // Check if next hop is a known blackhole
  let isBlackhole = blackholes.includes(nextHop);

  // Check if BGP community 65535:666 is set
  for (const c of communities) {
    if (c[0] === 65535 && c[1] === 666) {
      isBlackhole = true;
      break;
    }
  }

  if (isBlackhole) {
    return(
      <span className="route-prefix-flag blackhole-route is-blackhole-route">
        <FontAwesomeIcon icon={faCircle} />
        <div>Blackhole</div>
      </span>
    );
  }

  return (
    <span className="route-prefix-flag blackhole-route not-blackhole-route"></span>
  );
}

export default BlackholeIndicator;
