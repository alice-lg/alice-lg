
import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faSearch }
  from '@fortawesome/free-solid-svg-icons';


/**
 * The SearchInput is a text input field used for filtering
 */
const SearchInput = (props) => {
  return (
    <div className="input-group">
       <span className="input-group-addon">
        <FontAwesomeIcon icon={faSearch} />
       </span>
       <input type="text"
              className="form-control"
              {...props} />
    </div>
  );
};

export default SearchInput;
