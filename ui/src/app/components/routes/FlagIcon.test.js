import { render, screen }
  from '@testing-library/react';

import { faCircle }
  from '@fortawesome/free-solid-svg-icons';

import FlagIcon
  from './FlagIcon';

/**
 * Test rendering of the flag icon component.
 */
test('renders flag icon', () => {
  render(
    <div data-testid="icon">
      <FlagIcon icon={faCircle} tooltip="A flag icon" />
    </div>
  );

  // Check that the tooltip is in the document.
  expect(screen.getByText('A flag icon')).toBeInTheDocument();
});

