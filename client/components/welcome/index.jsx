
import React from 'react'

import PageHeader from 'components/page-header'

import Lookup from 'components/lookup'

export default class Welcome extends React.Component {
  render() {
    return (
      <div className="welcome-page">
       <PageHeader></PageHeader>

       <div className="jumbotron">
         <h1>Welcome to Birdseye!</h1>
         <p>Your friendly bird looking glass</p>
       </div>

			 <div className="col-md-8">
					<Lookup />
			 </div>

      </div>
    )
  }
}


