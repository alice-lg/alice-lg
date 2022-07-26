
import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';
import { faCircleExclamation }
  from '@fortawesome/free-solid-svg-icons';

import { useRejectCandidate }
  from 'app/context/bgp'; 


const RejectCandidateIndicator = ({route}) => {
  const isRejectCandidate = useRejectCandidate(route);
  if (!isRejectCandidate) {
    return null;
  }

  const cls = `route-prefix-flag reject-candidate-route`;
  return (
    <span className={cls}>
      <FontAwesomeIcon icon={faCircleExclamation} />
      <div>Reject Candidate</div>
    </span>
  );
}

export default RejectCandidateIndicator;
