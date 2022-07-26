
import { useEffect }
  from 'react';
import { useLocation } 
  from 'react-router-dom';


/**
 * ScrollToAnchor effect
 */
export const useScrollToAnchor = (refs) => {
  const { hash }  = useLocation();
  useEffect(() => {
    const ref = refs[hash];
    if (ref?.current) {
      ref.current.scrollIntoView();
    }
  }, [refs, hash]);
}

