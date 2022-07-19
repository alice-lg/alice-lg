
import { useConfig }
  from 'app/context/config';

import { resolveCommunities }
  from 'app/components/routes/communities'


const NoExportReason = ({route}) => {
  const config = useConfig();
  const noexportReasons = config.noexport_reasons;
  const routeCommunities = route?.bgp?.large_communities;

  if (!noexportReasons || !routeCommunities) {
      return null;
  }

  const reasons = resolveCommunities(
    noexportReasons, routeCommunities 
  );

  const reasonsView = reasons.map(([community, reason], key) => {
    const cls = `noexport-reason noexport-reason-${community[1]}-${community[2]}`;
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

export default NoExportReason;
