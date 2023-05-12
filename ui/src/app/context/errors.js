/**
 * Provide an error handler and an error state.
 */

import { useState
       , createContext
       , useContext
       , useCallback
       }
  from 'react';


const ErrorContext = createContext(null);
export const useErrors = () => useContext(ErrorContext);

export const useErrorHandler = () => {
  const [handle] = useErrors();
  return handle;
};


// Unfortunatley this does not really act as an error
// boundary. But we need to catch http errors from axios.
// Those are not cought using the ErrorBoundary approach.
export const ErrorsProvider = ({children}) => {
  const [errors, setErrors] = useState([]);

  // Handle prepends the error to the state.
  const handle = useCallback((err) => setErrors(
    (errors) => ([err, ...errors])), []);

  const dismiss = useCallback((err) => setErrors(
    (errors) => errors.filter((e) => e !== err)), []);

  const ctx = [handle, dismiss, errors];
  return (
    <ErrorContext.Provider value={ctx}>
      {children}
    </ErrorContext.Provider>
  );
}

// Check if the error (if present) has a status code
// that matches a gateway or a request timeout.
export const isTimeoutError = (error) => {
  const status = error?.response?.status;
  return (status === 504 || status === 408);
}
