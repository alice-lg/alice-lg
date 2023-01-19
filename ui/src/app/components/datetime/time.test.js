
import {parseServerTime} from 'app/components/datetime/time';

test("parse server time", () => {
  const t = "2023-10-24T23:42:11.3333333333Z";
  const result = parseServerTime(t).utc();
  expect(result).not.toBe(null);

  expect(result.month()).toBe(9);
  expect(result.year()).toBe(2023);
  expect(result.date()).toBe(24);
});
