import { useState
       , createContext
       , useEffect
       , useContext
       }
  from 'react';

import { updateContentApi } from 'api';

export const ContentContext = createContext({});

export const useContent = () => useContext(ContentContext);

export const ContentProvider = ({children}) => {
  const [content, setContent] = useState({});

  useEffect(() => {
    // Expose setContent in API
    updateContentApi(setContent);
  }, [setContent]);

  return (
    <ContentContext.Provider value={content}>
      {children}
    </ContentContext.Provider>
  );
};

