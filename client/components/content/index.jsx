
import react from 'react'
import {connect} from 'react-redux'

/*
 * Content Component
 */
function contentComponent(props) {
  if (!props.key) {
    return null;
  }

  // Traverse content by key, if content is found
  // return content, otherwise fall back to the default
  let tokens = key.split(".");
  let resolved = props.content;
  for (let part of tokens) {
    resolved = resolved[part];
    if (!resolved) {
      break;
    }
  }

  if (!resolved) {
    resolved = defaultValue; 
  }

  return (<span>{resolved}</span>);
}

export default connect(
  (state) => ({
    content: state.content
  })
)(contentComponent);

