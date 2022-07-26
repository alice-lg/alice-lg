
import { useMemo }
  from 'react';

import { useReadableCommunity }
  from 'app/context/bgp';


/*
 * Make style tags
 * Derive classes from community parts.
 */
const useStyeTags = (community) =>
  useMemo(() => (community.map(
    (part, i) =>
      `label-bgp-community-${i}-${part}`
    )).join(" "), [community]);

/*
 * Render community label
 */
const BgpCommunityLabel = ({community}) => {
  const readableCommunity = useReadableCommunity(community);
  const styleTags = useStyeTags(community);

  let cls = 'label label-bgp-community ';
  const label = community.join(":");

  if (!readableCommunity) {
    cls += "label-bgp-unknown";
    // Default label
    return (
      <span className={cls}>{label}</span>
    );
  }

  // Apply style
  cls += "label-info " + styleTags;
  return (<span className={cls}>{readableCommunity} ({label})</span>);
}

export default BgpCommunityLabel;
