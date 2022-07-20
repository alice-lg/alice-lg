
import GridLoader
  from 'react-spinners/GridLoader';

/**
 * Render a loading indicator that will
 * be visible if show is true
 */
const LoadingIndicator = ({show=true}) => {
  if (!show) {
    return null;
  }
  return (
    <div className="loading-indicator">
      <GridLoader loading={true} size={10} />
    </div>
  );
}

export default LoadingIndicator;
