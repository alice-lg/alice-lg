
import AsnLink
  from 'app/components/asns/AsnLink';

/**
 * Render an AS path as a list of links to ASNs.
 */
const AsPath = ({ asns }) => asns.map((asn, i) => (
  [<AsnLink key={i} asn={asn} />, " "]
));

export default AsPath;
