
import { Link } 
  from 'react-router-dom';

import { useNeighbors }
  from 'app/components/neighbors/Provider';


/**
 * Render Neighbors QuickLinks
 */
const QuickLinks = () => {
  const {isLoading} = useNeighbors();

  if (isLoading) {
    return null;
  }

  return (
    <div className="quick-links neighbors-quick-links">
      <span>Go to:</span>
      <ul>
        <li className="established">
          <Link to={{hash: "sessions-up"}}>Established</Link>
        </li>
        <li className="down">
          <Link to={{hash: "sessions-down"}}>Down</Link>
        </li>
      </ul>
    </div>
  );
}

export default QuickLinks;
