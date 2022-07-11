
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
          <a href="#sessions-up">Established</a>
        </li>
        <li className="down">
          <a href="#sessions-down">Down</a>
        </li>
      </ul>
    </div>
  );
}

export default QuickLinks;
