
import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faStar }
  from '@fortawesome/free-solid-svg-icons';

const PrimaryIndicator = ({route}) => {
  if (route.primary) {
    return(
      <span className="route-prefix-flag primary-route is-primary-route">
        <FontAwesomeIcon icon={faStar} />
        <div>Best Route</div>
      </span>
    );
  }
  return (
    <span className="route-prefix-flag primary-route not-primary-route"></span>
  );
}

export default PrimaryIndicator;
