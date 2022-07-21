
import { useConfig }
  from 'app/context/config';
import { useResolvedCommunities }
  from 'app/context/bgp'


const NoExportReason = ({route}) => {
  const { noexport_reasons } = useConfig();
  const communities = route?.bgp?.large_communities;
  const reasons = useResolvedCommunities(noexport_reasons, communities);

  if (!reasons) {
    return null;
  }

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
