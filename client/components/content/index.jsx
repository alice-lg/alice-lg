
import React from 'react'
import {connect} from 'react-redux'

/*
 * Content Component
 */
function ContentComponent(props) {
  let key = props.id;
  let defaultValue = props.children;

  if (!key) {
    return <span>{defaultValue}</span>;
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

  return (<span dangerouslySetInnerHTML={{__html: resolved}}></span>);
}

export default connect(
  (state) => ({
    content: state.content
  })
)(ContentComponent);

