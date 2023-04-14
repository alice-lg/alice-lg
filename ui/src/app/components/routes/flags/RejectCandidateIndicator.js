
import { faCircleExclamation }
  from '@fortawesome/free-solid-svg-icons';

import { useRejectCandidate }
  from 'app/context/bgp'; 

import FlagIcon
  from 'app/components/routes/FlagIcon';

/**
 * RejectCandidateIndicator
 * Displays a flag if the route is a reject candidate.
 *
 * @param route - The route to check
 */
const RejectCandidateIndicator = ({route}) => {
  const isRejectCandidate = useRejectCandidate(route);
  if (!isRejectCandidate) {
    return null;
  }

  return (
    <span className="route-prefix-flag reject-candidate-route">
      <FlagIcon icon={faCircleExclamation} tooltip="Reject Candidate" />
    </span>
  );
}

export default RejectCandidateIndicator;
