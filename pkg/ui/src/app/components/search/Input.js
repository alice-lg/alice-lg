
import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faSearch }
  from '@fortawesome/free-solid-svg-icons';

import { useEffect
       , useState
       }
  from 'react';

/**
 * The SearchInput is a text input field used for filtering.
 * The input is debounced and the onChange handler is called
 * with a delay.
 */
const SearchInput = ({value, onChange, debounce=0, ...props}) => {
  const [state, setState] = useState(value);

  useEffect(() => {
    const tRef = setTimeout(() => {
      onChange(state);
    }, debounce);
    return () => {
      clearTimeout(tRef); 
    };
  }, [state, debounce, onChange]);

  return (
    <div className="input-group">
       <span className="input-group-addon">
        <FontAwesomeIcon icon={faSearch} />
       </span>
       <input 
          type="text" className="form-control"
          onChange={(e) => setState(e.target.value)}
          {...props} />
    </div>
  );
};

export default SearchInput;
