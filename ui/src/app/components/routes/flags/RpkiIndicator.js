
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
  for (const com of communities) {
    // RPKI VALID
    if (com[0].toFixed() === rpkiValid[0] &&
        com[1].toFixed() === rpkiValid[1] &&
        com[2].toFixed() === rpkiValid[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-valid">
          <FlagIcon icon={faCircleCheck} tooltip="RPKI Valid" />
        </span>
      );
    }

    // RPKI UNKNOWN
    if (com[0].toFixed() === rpkiUnknown[0] &&
        com[1].toFixed() === rpkiUnknown[1] &&
        com[2].toFixed() === rpkiUnknown[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-unknown">
          <FlagIcon icon={faCircleQuestion} tooltip="RPKI Unknown" />
        </span>
      );
    }

    // RPKI NOT CHECKED
    if (com[0].toFixed() === rpkiNotChecked[0] &&
        com[1].toFixed() === rpkiNotChecked[1] &&
        com[2].toFixed() === rpkiNotChecked[2]) {
      return (
        <span className="route-prefix-flag rpki-route rpki-not-checked">
          <FlagIcon icon={faCircle} tooltip="RPKI Not Checked" />
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
          <FlagIcon icon={faCircleMinus} tooltip="RPKI Invalid" />
        </span>
      );
    }
  }

  return null;
}

export default RpkiIndicator;
