
const PrimaryIndicator = ({route}) => {
  if (route.primary) {
    return(
      <span className="route-prefix-flag primary-route is-primary-route"><i className="fa fa-star"></i>
        <div>Best Route</div>
      </span>
    );
  }
  return (
    <span className="route-prefix-flag primary-route not-primary-route"></span>
  );
}

export default PrimaryIndicator;
