
import { faStar }
  from '@fortawesome/free-solid-svg-icons';

import FlagIcon 
  from 'app/components/routes/FlagIcon';

/**
 * Show a primary route indicator icon
 *
 * @param route - The route object
 */
const PrimaryIndicator = ({route}) => {
  if (route.primary) {
    return(
      <span className="route-prefix-flag primary-route is-primary-route">
        <FlagIcon icon={faStar} tooltip="Best Route" />
      </span>
    );
  }
  return (
    <span className="route-prefix-flag primary-route not-primary-route"></span>
  );
}

export default PrimaryIndicator;
