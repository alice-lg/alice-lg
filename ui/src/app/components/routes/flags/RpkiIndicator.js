
import { useConfig }
  from 'app/context/config';

const RpkiIndicator = ({route}) => {
  const { rpki } = useConfig();

  // Check if indicator is enabled
  if (rpki.enabled === false) { return null; }

  // Check for BGP large communities as configured in the alice.conf
  const rpkiValid      = rpki.valid;
  const rpkiUnknown    = rpki.unknown;
  const rpkiNotChecked = rpki.not_checked;
  const rpkiInvalid    = rpki.invalid;

  const communities = route?.bgp?.large_communities;
  for (const com of communities) {
    // RPKI VALID
    if (com[0].toFixed() === rpkiValid[0] &&
        com[1].toFixed() === rpkiValid[1] &&
        com[2].toFixed() === rpkiValid[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-valid">
          <i className="fa fa-check-circle" /> 
          <div>RPKI Valid</div>
        </span>
      );
    }

    // RPKI UNKNOWN
    if (com[0].toFixed() === rpkiUnknown[0] &&
        com[1].toFixed() === rpkiUnknown[1] &&
        com[2].toFixed() === rpkiUnknown[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-unknown">
          <i className="fa fa-question-circle" />
          <div>RPKI Unknown</div>
        </span>
      );
    }

    // RPKI NOT CHECKED
    if (com[0].toFixed() === rpkiNotChecked[0] &&
        com[1].toFixed() === rpkiNotChecked[1] &&
        com[2].toFixed() === rpkiNotChecked[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-not-checked">
          <i className="fa fa-circle-o" />
          <div>RPKI not checked</div>
        </span>
      );
    }

    // RPKI INVALID
    // Depending on the configration this can either be a
    // single flag or a range with a given reason.
    let rpkiInvalidReason = 0;
    if (com[0].toFixed() === rpkiInvalid[0] &&
        com[1].toFixed() === rpkiInvalid[1]) {

      // This needs to be considered invalid, now try to detect why
      if (rpkiInvalid.length > 3 && rpkiInvalid[3] === "*") {
        // Check if token falls within range
        const start = parseInt(rpkiInvalid[2], 10);
        if (com[2] >= start) {
          rpkiInvalidReason = com[2];
        }
      } else {
        if (com[2].toFixed() === rpkiInvalid[2]) {
          rpkiInvalidReason = 1;
        }
      }
    }

    // This in invalid, render it!
    if (rpkiInvalidReason > 0) {
      const cls = `route-prefix-flag rpki-route rpki-invalid rpki-invalid-${rpkiInvalidReason}`
      return (
        <span className={cls}>
          <i className="fa fa-minus-circle" />
          <div>RPKI Invalid</div>
        </span>
      );
    }
  }

  return null;
}

export default RpkiIndicator;
