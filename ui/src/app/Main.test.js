
import { render } from '@testing-library/react';

import Main from 'app/Main';

test('render Main without crashing', () => {
  render(<Main />);
});
