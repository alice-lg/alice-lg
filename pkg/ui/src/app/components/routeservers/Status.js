
/**
 * RouteServer Status renders some information about
 * last reboot, etc. Also cache state will be displayed,
 * if provided.
 */

import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faClock
       , faThumbsUp
       , faThumbsDown
       }
  from '@fortawesome/free-regular-svg-icons';
import { faArrowsRotate }
  from '@fortawesome/free-solid-svg-icons';

import { useSelectedRouteServer }
  from 'app/components/routeservers/Provider';
import { useRouteServerStatus }
  from 'app/components/routeservers/StatusProvider';

import { useApiStatus }
  from 'app/components/api/StatusProvider';

import DateTime
  from 'app/components/datetime/DateTime';
import RelativeTime
  from 'app/components/datetime/RelativeTime';


/**
 * CacheStatus renders the current api cache status
 * from the context.
 */
const CacheStatus = () => {
  const status = useApiStatus();
  if (!status) {
    return null;
  }
  return (
   <tr>
     <td><FontAwesomeIcon icon={faArrowsRotate} /></td>
     <td>
       Generated <b><RelativeTime value={status.generatedAt}
                                  fuzzyNow={5}
                                  pastEvent={true} /></b>.<br />
       Next refresh <b><RelativeTime futureEvent={true}
                                     fuzzyNow={5}
                                     value={status.ttl} /></b>.
     </td>
   </tr>
  );
}

const Status = () => {
  const routeServer = useSelectedRouteServer();
  const rsStatus    = useRouteServerStatus();

  let lastReboot = rsStatus.last_reboot;
  if (lastReboot === "0001-01-01T00:00:00Z") {
      lastReboot = null;
  }

  let lastReconfig = rsStatus.last_reconfig;

  // We have some capabilities: openbgpd does not support
  // last reboot or last reconfig
  if (routeServer.type === "openbgpd") {
    lastReboot = null;
    lastReconfig = null;
  }

  return (
    <table className="routeserver-status">
      <tbody>
      {lastReboot &&
        <tr>
          <td><FontAwesomeIcon icon={faClock} /></td>
          <td>Last Reboot: <b><DateTime value={lastReboot} /></b></td>
        </tr>}
      {lastReconfig &&
        <tr>
          <td><FontAwesomeIcon icon={faClock} /></td>
          <td>Last Reconfig: <b><DateTime value={lastReconfig} /></b></td>
        </tr>}

      {rsStatus.message &&
        <tr>
          <td><FontAwesomeIcon icon={faThumbsUp} /></td>
          <td><b>{rsStatus.message}</b></td>
        </tr>}

      {!rsStatus.message &&
        <tr>
          <td><FontAwesomeIcon icon={faThumbsDown} /></td>
          <td>Route server is not reachable.</td>
        </tr>}

        <CacheStatus />
      </tbody>
    </table>
  );
}

export default Status;
