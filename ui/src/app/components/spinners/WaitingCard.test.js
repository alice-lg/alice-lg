
import {render} from '@testing-library/react';

import WaitingCard from 'app/components/spinners/WaitingCard';

beforeEach(() => {
  jest.useFakeTimers()
});

afterEach(() => {
  jest.runOnlyPendingTimers()
  jest.useRealTimers()
});

test("render waiting card", async () => {
  render(<WaitingCard isLoading={true} />);
});
