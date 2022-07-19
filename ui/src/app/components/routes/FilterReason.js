
import { useConfig }
  from 'app/context/config';

import { resolveCommunities }
  from 'app/components/routes/communities';

const FilterReason = ({route}) => {
  const config = useConfig();
  const rejectReasons = config.reject_reasons;

  const routeCommunities = route?.bgp?.large_communities;

  if (!rejectReasons || !routeCommunities) {
      return null;
  }

  const reasons = resolveCommunities(
    rejectReasons, routeCommunities,
  );

  const reasonsView = reasons.map(([community, reason], key) => {
    const cls = `reject-reason reject-reason-${community[1]}-${community[2]}`;
    return (
      <p key={key} className={cls}>
        <a href={`https://irrexplorer.nlnog.net/prefix/${route.network}`}
           rel="noreferrer"
           target="_blank" >{reason}</a>
      </p>
    );
  });

  return (<div className="reject-reasons">{reasonsView}</div>);
}

export default FilterReason;
