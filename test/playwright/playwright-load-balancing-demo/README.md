# Playwright Load Balancing Demo Test Suite

This Playwright test suite is designed to demonstrate the need for "Smart Load Balancing" in Testkube.

## Purpose

Standard sharding distributes tests round-robin (1, 2, 3, 4, 5, 1, 2...), which can lead to uneven execution times if long-running tests happen to fall to the same worker.

## Test Distribution

This suite contains 50 test files (`LB-test-01.spec.js` through `LB-test-50.spec.js`) in the `tests/` directory.

### Sleep Logic

- **Long tests (10 seconds)**: Files ending in 1 or 6 (01, 06, 11, 16, 21, 26, 31, 36, 41, 46)
- **Fast tests (100ms)**: All other files (02, 03, 04, 05, 07, 08, 09, 10, 12, 13, 14, 15, 17, 18, 19, 20, 22, 23, 24, 25, 27, 28, 29, 30, 32, 33, 34, 35, 37, 38, 39, 40, 42, 43, 44, 45, 47, 48, 49, 50)

## Expected Behavior with 5 Shards (Round-Robin)

When executed with 5 parallel workers using standard round-robin distribution:

- **Worker 1**: Gets tests 1, 6, 11, 16, 21, 26, 31, 36, 41, 46 (all 10s each) = **~100 seconds total**
- **Worker 2**: Gets tests 2, 7, 12, 17, 22, 27, 32, 37, 42, 47 (all 100ms each) = **~1 second total**
- **Worker 3**: Gets tests 3, 8, 13, 18, 23, 28, 33, 38, 43, 48 (all 100ms each) = **~1 second total**
- **Worker 4**: Gets tests 4, 9, 14, 19, 24, 29, 34, 39, 44, 49 (all 100ms each) = **~1 second total**
- **Worker 5**: Gets tests 5, 10, 15, 20, 25, 30, 35, 40, 45, 50 (all 100ms each) = **~1 second total**

This demonstrates the problem: one worker takes 100 seconds while others finish in 1 second, proving the need for smart load balancing that considers test execution time.

## Shard Files

The `shards/` directory contains `shard-1.txt` through `shard-5.txt`, each listing the test file paths (one per line) for that shard in round-robin order.

## Regenerating Files

To regenerate all test files:

```bash
node generate-tests.js
```

To regenerate shard files:

```bash
node generate-shards.js
```

## Running the Tests

This is a standard Playwright project. Install dependencies and run tests as usual:

```bash
npm install
npx playwright install
npx playwright test
```

To run only the tests from a specific shard file (e.g. shard 1):

```bash
npx playwright test $(cat shards/shard-1.txt | tr '\n' ' ')
```
