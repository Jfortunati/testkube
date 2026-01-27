const fs = require('fs');
const path = require('path');

const shardsDir = path.join(__dirname, 'cypress', 'shards');

// Ensure directory exists
if (!fs.existsSync(shardsDir)) {
  fs.mkdirSync(shardsDir, { recursive: true });
}

// Generate 5 shard files with round-robin distribution
const numShards = 5;
const numTests = 50;

// Initialize shard arrays
const shards = Array.from({ length: numShards }, () => []);

// Distribute tests in round-robin fashion
for (let i = 1; i <= numTests; i++) {
  const testNumber = i.toString().padStart(2, '0');
  const testFile = `cypress/e2e/LB-test-${testNumber}.spec.js`;
  const shardIndex = (i - 1) % numShards;
  shards[shardIndex].push(testFile);
}

// Write each shard file
for (let shardIndex = 0; shardIndex < numShards; shardIndex++) {
  const shardNumber = shardIndex + 1;
  const fileName = `shard-${shardNumber}.txt`;
  const filePath = path.join(shardsDir, fileName);
  
  // Write one test file per line
  const content = shards[shardIndex].join('\n') + '\n';
  
  fs.writeFileSync(filePath, content, 'utf8');
  console.log(`Generated ${fileName} with ${shards[shardIndex].length} test files`);
  console.log(`  Tests: ${shards[shardIndex].map(f => f.match(/LB-test-(\d+)/)[1]).join(', ')}`);
}

console.log('\nAll shard files generated successfully!');
console.log('\nShard distribution:');
console.log('Shard 1: Gets tests ending in 1 or 6 (10s each) = ~100 seconds');
console.log('Shards 2-5: Get other tests (100ms each) = ~1 second each');

