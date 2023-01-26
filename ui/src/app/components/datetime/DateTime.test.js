import {render, screen} from '@testing-library/react';

import DateTime from 'app/components/datetime/DateTime';


test("render a parsed server time as date time", () => {
  const t = "2022-05-06T23:42:11.123Z";
  render(<p data-testid="result"><DateTime value={t} utc={true}/></p>);

  const result = screen.getByTestId("result");
  expect(result.innerHTML).toBe("Friday, May 6, 2022 11:42 PM");
});
