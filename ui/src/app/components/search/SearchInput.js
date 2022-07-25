
import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faSearch }
  from '@fortawesome/free-solid-svg-icons';

import { useEffect
       , useState
       , useRef
       , forwardRef
       }
  from 'react';

/**
 * The SearchInput is a text input field used for filtering.
 * The input is debounced and the onChange handler is called
 * with a delay.
 */
const SearchInput = forwardRef(({
  value,
  onChange,
  debounce=0,
  ...props
}, ref) => {
  const valueRef = useRef();
  const [state, setState] = useState(value);

  useEffect(() => {
    // Prevent update loops and unwanted calls
    if (valueRef.current === undefined && state === "") {
      return;
    }
    if (valueRef.current === state) {
      return;
    }
    const tRef = setTimeout(() => {
      onChange(state);
    }, debounce);

    valueRef.current = state;
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
          value={state}
          onChange={(e) => setState(e.target.value)}
          ref={ref}
          {...props} />
    </div>
  );
});

export default SearchInput;
