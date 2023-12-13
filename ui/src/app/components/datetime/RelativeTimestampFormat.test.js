import {render, screen} from '@testing-library/react';
import moment from 'moment';

import RelativeTimestampFormat from 'app/components/datetime/RelativeTimestampFormat';


test("render a formatted relative timestamp", () => {
  const now = moment.utc();
  const time = now.clone().subtract(10, 'hours');
  const t = (now - time) * 1000 * 1000;
  render(<p data-testid="result"><RelativeTimestampFormat value={t} /></p>);
  const result = screen.getByTestId("result");

  const expected = time.format();
  expect(result.innerHTML).toBe(expected);
});


