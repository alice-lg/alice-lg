
import { useEffect
       , useCallback
       }
  from 'react';
import { useNavigate
       , useLocation
       }
  from 'react-router-dom';

import { useConfig }
  from 'app/context/config';
import { useQuery }
  from 'app/context/query';

import PageHeader
  from 'app/components/page/Header';
import Content
  from 'app/components/content/Content';
import SearchGlobalInput
  from 'app/components/search/SearchGlobalInput';


const StartGlobalSearch = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { prefix_lookup_enabled } = useConfig();
  const [{q}] = useQuery({q: ""});

  const openSearch = useCallback((qry) => {
    if (!prefix_lookup_enabled || qry === "") {
      return;
    }
    navigate({...location, pathname: "/search"}, {replace: true});
  }, [prefix_lookup_enabled, navigate, location]);

  // Make sure location is search when we are running a query
  useEffect(() => {
    openSearch(q);
  }, [q, openSearch]);

  if (!prefix_lookup_enabled) {
    return null;
  }

  return (
    <div className="lookup-container">
     <div className="col-md-8">
       <SearchGlobalInput />
     </div>
    </div>
  );
}


const StartPage = () => {
  return (
    <div className="welcome-page">
      <PageHeader></PageHeader>

      <div className="jumbotron">
        <h1><Content id="welcome.title">Welcome to Alice!</Content></h1>
        <p><Content id="welcome.tagline">Your friendly BGP looking glass</Content></p>
      </div>

      <StartGlobalSearch />
    </div>
  );
}

export default StartPage;
