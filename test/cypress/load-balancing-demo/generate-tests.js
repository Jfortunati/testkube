const fs = require('fs');
const path = require('path');

const e2eDir = path.join(__dirname, 'cypress', 'e2e');

// Ensure directory exists
if (!fs.existsSync(e2eDir)) {
  fs.mkdirSync(e2eDir, { recursive: true });
}

// Generate 50 test files
for (let i = 1; i <= 50; i++) {
  const testNumber = i.toString().padStart(2, '0');
  const fileName = `LB-test-${testNumber}.spec.js`;
  const filePath = path.join(e2eDir, fileName);
  
  // Determine sleep duration: files ending in 1 or 6 (index % 5 == 1) sleep for 10 seconds
  // All others sleep for 100ms
  const lastDigit = i % 10;
  const isLongTest = (lastDigit === 1 || lastDigit === 6);
  const sleepDuration = isLongTest ? 10000 : 100;
  
  const comment = isLongTest 
    ? `// This test sleeps for 10 seconds (files ending in 1 or 6: 01, 06, 11, 16, 21, 26, 31, 36, 41, 46)`
    : `// This test sleeps for 100ms (all other files)`;
  
  const testContent = `describe('Load Balancing Test ${testNumber}', () => {
  it('simulated duration', async () => {
    ${comment}
    // When run with 5 shards using round-robin, Worker 1 gets tests 1,6,11,16,21,26,31,36,41,46 (all 10s each = 100s total)
    // Other workers get fast tests (100ms each = ~1s total), demonstrating the need for smart load balancing
    await new Promise(resolve => setTimeout(resolve, ${sleepDuration}));
    expect(true).to.be.true;
  });
});
`;

  fs.writeFileSync(filePath, testContent, 'utf8');
  console.log(`Generated ${fileName}`);
}

console.log('\nAll 50 test files generated successfully!');

