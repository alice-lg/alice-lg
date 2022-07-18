import { useState
       , createContext
       , useEffect
       , useContext
       }
  from 'react';

export const ContentContext = createContext({});

export const useContent = () => useContext(ContentContext);

export const ContentProvider = ({children}) => {
  const [content, setContent] = useState({});

  // Expose setContent as API??
  useEffect(() => {
    if (!window.API) {
      window.API = {};
    }
    window.API.setContent = setContent;
  }, [setContent]);

  return (
    <ContentContext.Provider value={content}>
      {children}
    </ContentContext.Provider>
  );
};

