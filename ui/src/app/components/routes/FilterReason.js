
import { useConfig }
  from 'app/context/config';
import { useResolvedCommunities }
  from 'app/context/bgp';


const FilterReason = ({route}) => {
  const { reject_reasons } = useConfig();
  const communities = route?.bgp?.large_communities;
  const reasons = useResolvedCommunities(reject_reasons, communities);

  if (!reasons) {
      return null;
  }

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
