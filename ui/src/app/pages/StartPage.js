

import PageHeader
  from 'app/components/page/Header';
import Content
  from 'app/components/content/Content';

const StartPage = () => {
  return (
    <div className="welcome-page">
     <PageHeader></PageHeader>

     <div className="jumbotron">
       <h1><Content id="welcome.title">Welcome to Alice!</Content></h1>
       <p><Content id="welcome.tagline">Your friendly bird looking glass</Content></p>
     </div>

     {/*<LookupWidget />*/}

    </div>
  );
}

export default StartPage;
