const { defineConfig } = require('@playwright/test');

// Fixed: Both 'list' and 'junit' are now defined as tuples inside the array
const reporters = process.env.PLAYWRIGHT_JUNIT_OUTPUT_NAME
  ? [
      ['list'], 
      ['junit', { outputFile: process.env.PLAYWRIGHT_JUNIT_OUTPUT_NAME }]
    ]
  : 'list';

module.exports = defineConfig({
  testDir: './tests',
  timeout: 60 * 1000,
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: 0,
  reporter: reporters, // This will now receive a valid nested array
  use: {
    baseURL: 'https://testkube-test-page-lipsum.pages.dev/',
    video: 'off',
    screenshot: 'off',
    trace: 'off',
  },
});
