const { test, expect } = require('@playwright/test');

test('Load Balancing Test 25 - simulated duration', async () => {
  // This test sleeps for 100ms (all other files)
  // When run with 5 shards using round-robin, Worker 1 gets tests 1,6,11,16,21,26,31,36,41,46 (all 10s each = 100s total)
  // Other workers get fast tests (100ms each = ~1s total), demonstrating the need for smart load balancing
  await new Promise((resolve) => setTimeout(resolve, 100));
  expect(true).toBe(true);
});
