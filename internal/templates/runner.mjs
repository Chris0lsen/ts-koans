import { readFileSync } from "fs";
import vm from "vm";

const combined = readFileSync("./run.js", "utf8");

// Optionally: set up a basic sandbox (exports/global, etc.)
const sandbox = { exports: {}, module: { exports: {}}, console };
const context = vm.createContext(sandbox);

try {
  vm.runInContext(combined, context, { timeout: 1000 });
  console.log("✅ All tests passed!")
} catch (err) {
  // Print clean error message
  console.log("❌ Test failed:", err && err.message ? err.message : err);
  process.exit(1);
}
