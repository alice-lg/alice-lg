import {render, screen} from '@testing-library/react';
import moment from 'moment';

import RelativeTime from 'app/components/datetime/RelativeTime';


test("render a relative time", () => {
  const t = moment().subtract(11, 'hours');
  render(<p data-testid="result"><RelativeTime value={t} /></p>);

  const result = screen.getByTestId("result");
  expect(result.innerHTML).toBe("11 hours ago");
});


