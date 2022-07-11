
/**
 * The SearchInput is a text input field used for filtering
 */
const SearchInput = (props) => {
  return (
    <div className="input-group">
       <span className="input-group-addon">
        <i className="fa fa-search"></i>
       </span>
       <input type="text"
              className="form-control"
              {...props} />
    </div>
  );
};

export default SearchInput;
