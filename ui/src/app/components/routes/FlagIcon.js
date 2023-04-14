import { FontAwesomeIcon }
  from '@fortawesome/react-fontawesome';

/**
 * Display a route flag icon indicator with a tooltip.
 * 
 * @param icon - The icon to display.
 * @param tooltip - The tooltip to display.
 */
const FlagIcon = ({icon, tooltip}) => {
  return (
    <>
      <i><FontAwesomeIcon icon={icon} /></i>
      <div>{tooltip}</div>
    </>
  );
}

export default FlagIcon;


