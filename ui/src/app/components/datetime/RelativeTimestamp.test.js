import {render, screen} from '@testing-library/react';

import RelativeTimestamp from 'app/components/datetime/RelativeTimestamp';


test("render a relative timestamp", () => {
  const t = 15 * 60 * 1000 * 1000 * 1000; // 15 min
  render(<p data-testid="result"><RelativeTimestamp value={t} /></p>);

  const result = screen.getByTestId("result");
  expect(result.innerHTML).toBe("15 minutes ago");
});


