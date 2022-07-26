
import { useContent }
  from 'app/context/content';

const Content = ({id, children}) => {
  let defaultValue = children;
  let content = useContent();

  if (!id) {
    return defaultValue;
  }

  // Traverse content by key, if content is found
  // return content, otherwise fall back to the default
  let tokens = id.split(".");
  let resolved = content;
  for (let part of tokens) {
    resolved = resolved[part];
    if (!resolved) {
      break;
    }
  }

  if (!resolved) {
    resolved = defaultValue; 
  }

  return <span dangerouslySetInnerHTML={{__html: resolved}}></span>;
};

export default Content;
