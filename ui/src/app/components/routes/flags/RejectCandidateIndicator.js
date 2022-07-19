
import { useConfig }
  from 'app/context/config';

import { isRejectCandidate }
  from 'app/components/routes/communities'; 


const RejectCandidateIndicator = ({route}) => {
  const { reject_candidates } = useConfig();
  const candidateCommunities = reject_candidates.communities;

  if (candidateCommunities) {
    return null;
  }
  if (!isRejectCandidate(candidateCommunities, route)) {
    return null;
  }

  const cls = `route-prefix-flag reject-candidate-route`;
  return (
    <span className={cls}>
      <i className="fa fa-exclamation-circle" />
      <div>Reject Candidate</div>
    </span>
  );
}

export default RejectCandidateIndicator;
