import { render, screen }
  from '@testing-library/react';

import { ConfigContext }
  from 'app/context/config';
import RpkiIndicator
  from './RpkiIndicator';

// Provide config context with rpki settings
const config = {
  rpki: {
    enabled: true,
    valid: [["1234", "1111", "1"]],
    unknown: [["1234", "1111", "0"]],
    not_checked: [["1234", "1111", "10"]],
    invalid: [["1234", "1111", "100"]],
  },
};


/**
 * Test rendering the RpkiIndicator component.
 */
test('renders RpkiIndicator with a valid route', () => {
  // Render the RpkiIndicator component for a valid prefix
  const route = {
    bgp: {
      large_communities: [
        [1234, 1111, 1],
      ],
    },
  };
  render(
    <ConfigContext.Provider value={config}>
      <RpkiIndicator route={route} />
    </ConfigContext.Provider>
  );
  expect(screen.getByText('RPKI Valid')).toBeInTheDocument();
});

/**
 * Test rendering the RpkiIndicator component with 
 * an rpki unknown route.
 */
test('renders RpkiIndicator with an unknown route', () => {
  // Render the RpkiIndicator component for an unknown prefix
  const route = {
    bgp: {
      large_communities: [
        [1234, 1111, 0],
      ],
    },
  };
  render(
    <ConfigContext.Provider value={config}>
      <RpkiIndicator route={route} />
    </ConfigContext.Provider>
  );
  expect(screen.getByText('RPKI Unknown')).toBeInTheDocument();
});

/**
 * Test rendering the RpkiIndicator component with
 * an rpki not checked route.
 */
test('renders RpkiIndicator with a not checked route', () => {
  // Render the RpkiIndicator component for a not checked prefix
  const route = {
    bgp: {
      large_communities: [
        [1234, 1111, 10],
      ],
    },
  };
  render(
    <ConfigContext.Provider value={config}>
      <RpkiIndicator route={route} />
    </ConfigContext.Provider>
  );
  expect(screen.getByText('RPKI Not Checked')).toBeInTheDocument();
});

/**
 * Test rendering the RpkiIndicator component with an
 * rpki invalid route.
 */
test('renders RpkiIndicator with an invalid route', () => {
  // Render the RpkiIndicator component for an invalid prefix
  const route = {
    bgp: {
      large_communities: [
        [1234, 1111, 100],
      ],
    },
  };
  render(
    <ConfigContext.Provider value={config}>
      <RpkiIndicator route={route} />
    </ConfigContext.Provider>
  );
  expect(screen.getByText('RPKI Invalid')).toBeInTheDocument();
});
