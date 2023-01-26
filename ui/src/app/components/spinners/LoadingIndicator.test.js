import {render} from '@testing-library/react';

import LoadingIndicator from 'app/components/spinners/LoadingIndicator';

test("render loading indicator", () => {
  render(<LoadingIndicator show={true} />);
});
