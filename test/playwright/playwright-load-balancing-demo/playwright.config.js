const { defineConfig } = require('@playwright/test');

module.exports = defineConfig({
  testDir: './tests',
  timeout: 60 * 1000,
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: 0,
  reporter: 'list',
  use: {
    baseURL: 'https://testkube-test-page-lipsum.pages.dev/',
    video: 'off',
    screenshot: 'off',
    trace: 'off',
  },
});
