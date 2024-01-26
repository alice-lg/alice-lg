
/**
 * Wrap an ASNumber with a link to bgp.tools for more information.
 */
const AsnLink = ({ asn }) => {
    // const baseUrl = "https://irrexplorer.nlnog.net/asn/AS";
    const baseUrl = "https://bgp.tools/as/";
    const url = `${baseUrl}${asn}`;
    return (
      <a href={url} target="_blank" rel="noreferrer">{asn}</a> 
    );
}

export default AsnLink;
