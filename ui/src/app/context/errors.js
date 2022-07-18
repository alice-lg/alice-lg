/**
 * Provide an error handler and an error state.
 */

import { useState
       , createContext
       , useContext
       , useRef
       , useCallback
       }
  from 'react';


const ErrorContext = createContext(null);
export const useErrors = () => useContext(ErrorContext);

export const useErrorHandler = () => {
  const [handleRef] = useErrors();
  return useCallback((err) => handleRef.current(err), [handleRef]);
};


// Unfortunatley this does not really act as an error
// boundary. But we need to catch http errors from axios.
// Those are not cought using the ErrorBoundary approach.
export const ErrorsProvider = ({children}) => {
  const [errors, setErrors] = useState([]);

  // Handle prepends the error to the state.
  // Use a ref to the handler function to prevent
  // a rendering loop.
  const handle = (err) => {
    setErrors([err, ...errors]);
  };
  const handleRef = useRef(handle);

  // Dismiss removes the error from the state
  const dismiss = (err) => {
    const filtered = errors.filter((e) => e !== err)
    setErrors(filtered);
  }

  const ctx = [handleRef, dismiss, errors];
  return (
    <ErrorContext.Provider value={ctx}>
      {children}
    </ErrorContext.Provider>
  );
}

