import { useState
       , createContext
       }
  from 'react';

export const ContentContext = createContext({});

const ContentProvider = ({children}) => {
  const [content, setContent] = useState({});

  // Expose setContent as API??

  return (
    <ContentContext.Provider value={content}>
      {children}
    </ContentContext.Provider>
  );
};

export default ContentProvider;
