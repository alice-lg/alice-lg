
/**
 * Errors Component renders a stack of dismissable errors.
 */

import { useRouteservers }
  from 'app/components/routeservers/Provider';

import { useErrors }
  from 'app/components/errors/Provider';


/**
 * infoFromError extracts error data
 */
const infoFromError = (error) => {
    if (error.response && error.response.data && error.response.data.code) {
      return error.response.data;
    }
    return null;
}


/**
 * Error renders a single dismissable error
 */
const Error = ({error, onDismiss}) => {
  const routeservers = useRouteservers();   

  let status = 600;
  if (error.response) {
    status = error.response.status;
  }
  if (!status || (status !== 429 && status < 500)) {
    return null;
  }

  const errorInfo = infoFromError(error);

  // Find affected routeserver
  let rs = null;
  if (errorInfo) {
    const rsId = errorInfo.routeserver_id; 
    if (rsId !== null) {
      rs = routeservers.find(r => r.id === rsId);
    }
  }

  let body = null;

  if (status == 429) {
    body = (
      <div className="error-message">
        <p>Alice reached the request limit.</p>
        <p>We suggest you try at a less busy time.</p>
      </div>
    );
  } else {
    let errorStatus = "";
    if (error.response) {
      errorStatus = " (got HTTP " + error.response.status + ")";
    }
    if (errorInfo) {
      errorStatus = ` (got ${errorInfo.tag})`;
    }

    body = (
      <div className="error-message">
        <p>Alice has trouble connecting to the API 
          {rs && <span> of <b>{rs.name}</b></span>}
          {errorStatus}.
        </p>
        <p>If this problem persist, we suggest you
        try again later.</p>
      </div>
    );
  }

  return (
    <div className="error-notify">
      <div className="error-dismiss">
        <i className="fa fa-times-circle" aria-hidden="true"
           onClick={() => onDismiss(error)}></i>
      </div>
      <div className="error-icon">
        <i className="fa fa-times-circle" aria-hidden="true"></i>
      </div>
      {body}
    </div>
  );
}

/**
 * Errors displays a stacked errors list
 */
const Errors = () => {
  const [_, dismiss, errors] = useErrors();
  return errors.map((err, i) =>
    <Error
      key={i}
      error={err}
      onDismiss={(err) => dismiss(err)} />);
}

export default Errors;
