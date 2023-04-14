import { render, screen }
  from '@testing-library/react';

import RejectCandidateIndicator
  from './RejectCandidateIndicator';

import { ConfigContext }
  from 'app/context/config';


// Mock config with reject candidate community
const config = {
  reject_candidates: {
    communities: {
      1111: {
        1234: {
          1: "reject-candidate-2",
        },
      },
    },
  },
};

/**
 * Test the RejectCandidateIndicator component with
 */
test('renders reject candidate indicator' , () => {
  const route = {
    bgp: {
      large_communities: [
        [1111, 1234, 1],
      ],
    },
  };

  // Render the component
  render(
    <ConfigContext.Provider value={config}>
      <RejectCandidateIndicator route={route} />
    </ConfigContext.Provider>
  );

  // Check that the indicator is rendered
  expect(screen.getByText('Reject Candidate')).toBeInTheDocument();
});

