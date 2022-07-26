
import { Link
       , useLocation
       } 
  from 'react-router-dom';

import { useNeighbors }
  from 'app/context/neighbors';


/**
 * Render Neighbors QuickLinks
 */
const QuickLinks = () => {
  const location = useLocation();
  const {isLoading} = useNeighbors();

  if (isLoading) {
    return null;
  }

  return (
    <div className="quick-links neighbors-quick-links">
      <span>Go to:</span>
      <ul>
        <li className="established">
          <Link to={{...location, hash: "sessions-up"}}>Established</Link>
        </li>
        <li className="down">
          <Link to={{...location, hash: "sessions-down"}}>Down</Link>
        </li>
      </ul>
    </div>
  );
}

export default QuickLinks;
