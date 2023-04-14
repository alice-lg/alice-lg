
import { render, screen }
  from '@testing-library/react';

import PrimaryIndicator
  from './PrimaryIndicator';

/**
 * Test rendering the primary indicator
 */
test('renders primary indicator', () => {
  // Routes for testing: primary and not primary
  const primaryRoute = {
    primary: true,
  };
  const notPrimaryRoute = {
    primary: false,
  };

  // Render the non primary route indicator
  render(<PrimaryIndicator route={notPrimaryRoute} />);
  expect(screen.queryByText('Best Route')).not.toBeInTheDocument();

  // Render the primary indicator
  render(<PrimaryIndicator route={primaryRoute} />);
  expect(screen.getByText('Best Route')).toBeInTheDocument();
});

