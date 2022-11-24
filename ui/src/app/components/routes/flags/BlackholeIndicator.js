
import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faCircle }
  from '@fortawesome/free-solid-svg-icons';

import { useRouteServer }
  from 'app/context/route-servers';
import { matchCommunityRange
       , useBlackholeCommunities 
       }
  from 'app/context/bgp';


const BlackholeIndicator = ({route}) => {
  const routeServer = useRouteServer(); // blackholes are store per RS
  const blackholeCommunities = useBlackholeCommunities();

  const blackholes = routeServer?.blackholes || [];
  const nextHop = route?.bgp?.next_hop;
  const routeStandard = route?.bgp?.communities || [];
  const routeExtended = route?.bgp?.ext_communities || [];
  const routeLarge    = route?.bgp?.large_communities || [];

  // Check if next hop is a known blackhole
  let isBlackhole = blackholes.includes(nextHop);

  // Check standard communities
  for (const c of blackholeCommunities.standard) {
    for (const r of routeStandard) {
      if (matchCommunityRange(r, c)) {
        isBlackhole = true;
        break;
      }
    }
  }
  // Check large communities
  for (const c of blackholeCommunities.large) {
    for (const r of routeLarge) {
      if (matchCommunityRange(r, c)) {
        isBlackhole = true;
        break;
      }
    }
  }
  // Check extended
  for (const c of blackholeCommunities.extended) {
    for (const r of routeExtended) {
      if (matchCommunityRange(r, c)) {
        isBlackhole = true;
        break;
      }
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
