describe('Load Balancing Test 06', () => {
  it('simulated duration', async () => {
    // This test sleeps for 10 seconds (files ending in 1 or 6: 01, 06, 11, 16, 21, 26, 31, 36, 41, 46)
    // When run with 5 shards using round-robin, Worker 1 gets tests 1,6,11,16,21,26,31,36,41,46 (all 10s each = 100s total)
    // Other workers get fast tests (100ms each = ~1s total), demonstrating the need for smart load balancing
    await new Promise(resolve => setTimeout(resolve, 10000));
    expect(true).to.be.true;
  });
});
