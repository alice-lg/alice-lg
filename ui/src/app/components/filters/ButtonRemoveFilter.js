
import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faTimes }
  from '@fortawesome/free-solid-svg-icons';

const ButtonRemoveFilter = ({onClick}) => (
  <button
    className="btn btn-remove"
    onClick={onClick}>
      <FontAwesomeIcon icon={faTimes} />
  </button>
);

export default ButtonRemoveFilter;
