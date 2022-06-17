/**
 * Provide an error handler and an error state.
 */

import { Component
       , useState
       , createContext
       , useContext
       }
  from 'react';


const ErrorContext = createContext(null);
export const useErrors = () => useContext(ErrorContext);


// Unfortunatley this does not really act as an error
// boundary. But we need to catch http errors from axios.
// Those are not cought using the ErrorBoundary approach.
const ErrorProvider = ({children}) => {
  const [errors, setErrors] = useState([]);

  const handle = (err) => {
    setErrors([err, ...errors]);
  }

  const dismiss = (err) => {
    const filtered = errors.filter((e) => e != err)
    setErrors(filtered);
  }

  const ctx = [handle, dismiss, errors];

  return (
    <ErrorContext.Provider value={ctx}>
      {children}
    </ErrorContext.Provider>
  );
}

export default ErrorProvider;
