
import { faCircleCheck
       , faCircleMinus
       , faCircleQuestion
       }
  from '@fortawesome/free-solid-svg-icons';
import { faCircle }
  from '@fortawesome/free-regular-svg-icons';

import { useConfig }
  from 'app/context/config';

import FlagIcon
  from 'app/components/routes/FlagIcon';

const RpkiIndicator = ({route}) => {
  const { rpki } = useConfig();

  // Check if indicator is enabled
  if (rpki.enabled === false) { return null; }

  // Check for BGP large communities as configured in the alice.conf
  // FIXME: why are we using strings here?! ['1234', '123', '1'].
  const rpkiValid      = rpki.valid;
  const rpkiUnknown    = rpki.unknown;
  const rpkiNotChecked = rpki.not_checked;
  const rpkiInvalid    = rpki.invalid;

  const communities = route?.bgp?.large_communities || [];

  const matchCommunity = (com, coms) =>
    coms.some((match) =>
       (com[0].toFixed() === match[0] &&
        com[1].toFixed() === match[1] &&
        com[2].toFixed() === match[2]));

  for (const com of communities) {
    // RPKI VALID
    if (matchCommunity(com, rpkiValid)) {
      return (
        <span className="route-prefix-flag rpki-route rpki-valid">
          <FlagIcon icon={faCircleCheck} tooltip="RPKI Valid" />
        </span>
      );
    }

    // RPKI UNKNOWN
    if (matchCommunity(com, rpkiUnknown)) {
      return (
        <span className="route-prefix-flag rpki-route rpki-unknown">
          <FlagIcon icon={faCircleQuestion} tooltip="RPKI Unknown" />
        </span>
      );
    }

    // RPKI NOT CHECKED
    if (matchCommunity(com, rpkiNotChecked)) {
      return (
        <span className="route-prefix-flag rpki-route rpki-not-checked">
          <FlagIcon icon={faCircle} tooltip="RPKI Not Checked" />
        </span>
      );
    }

    // RPKI INVALID
    // Depending on the configuration this can either be a
    // single flag or a range with a given reason.
    let rpkiInvalidReason = 0;
    for (const invalid of rpkiInvalid) {
      if (com[0].toFixed() === invalid[0] &&
          com[1].toFixed() === invalid[1]) {

        // This needs to be considered invalid, now try to detect why
        if (invalid.length > 3 && invalid[3] === "*") {
          // Check if token falls within range
          const start = parseInt(invalid[2], 10);
          if (com[2] >= start) {
            rpkiInvalidReason = com[2];
          }
        } else {
          if (com[2].toFixed() === invalid[2]) {
            rpkiInvalidReason = 1;
          }
        }
        break; // We found a match, stop searching
      }
    }

    // This in invalid, render it!
    if (rpkiInvalidReason > 0) {
      const cls = `route-prefix-flag rpki-route rpki-invalid rpki-invalid-${rpkiInvalidReason}`
      return (
        <span className={cls}>
          <FlagIcon icon={faCircleMinus} tooltip="RPKI Invalid" />
        </span>
      );
    }
  }

  return null;
}

export default RpkiIndicator;
