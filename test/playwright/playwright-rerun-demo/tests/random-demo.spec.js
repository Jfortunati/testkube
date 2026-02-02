// @ts-check
const { test, expect } = require('@playwright/test');

const PASS_CHANCE = 0.7;

function shouldPass() {
  return Math.random() < PASS_CHANCE;
}

for (let i = 1; i <= 50; i++) {
  test(`random-demo-${i}`, async () => {
    if (shouldPass()) {
      expect(1).toBe(1);
    } else {
      expect(1).toBe(2);
    }
  });
}
