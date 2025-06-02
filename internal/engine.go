package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func copyVersionFilesToTempDir(tempDir string) {
	for _, filename := range []string{".tool-versions", ".nvmrc", ".node-version"} {
		if data, err := os.ReadFile(filename); err == nil {
			os.WriteFile(filepath.Join(tempDir, filename), data, 0600)
		}
	}
}

func RunTypeScript(userCode, testScript string) (string, error) {
	dir, err := os.MkdirTemp("", "tskoans")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	copyVersionFilesToTempDir(dir)
	defer os.RemoveAll(dir)

	userTs := filepath.Join(dir, "user.ts")
	runnerJs := filepath.Join(dir, "runner.mjs")

	// Write user TypeScript code
	if err := os.WriteFile(userTs, []byte(userCode), 0600); err != nil {
		return "", fmt.Errorf("failed to write user code: %w", err)
	}

	// Compile TS to JS
	tscCmd := exec.Command("tsc", userTs, "--target", "es2020", "--module", "commonjs", "--outDir", dir)
	tscCmd.Dir = dir
	tscOut, tscErr := tscCmd.CombinedOutput()
	if tscErr != nil {
		return string(tscOut), nil // Return TypeScript compilation errors
	}

	// Write runner.mjs
	runner := fmt.Sprintf(`
import { readFileSync } from 'fs';
import vm from 'vm';

console.log("[runner] Starting test runner...");

let userCode;
try {
    userCode = readFileSync('./user.js', 'utf8');
    console.log("[runner] Loaded user.js");
} catch (e) {
    console.log("[runner] Failed to load user.js:", e.message);
    process.exit(1);
}

const sandbox = {
    exports: {},
    module: { exports: {} },
    console,
};

vm.createContext(sandbox);

try {
    vm.runInContext(userCode, sandbox, { timeout: 1000 });
    const exports = sandbox.exports;
    // Your test script uses 'exports' as in 'exports.double'
    %s
    console.log("✅ All tests passed!");
} catch (e) {
    console.log("❌ Test failed:", e && e.message ? e.message : e);
    process.exit(1);
}
`, testScript)

	if err := os.WriteFile(runnerJs, []byte(runner), 0600); err != nil {
		return "", fmt.Errorf("failed to write runner: %w", err)
	}

	// Run Node on the runner
	nodeCmd := exec.Command("node", runnerJs)
	nodeCmd.Dir = dir
	out, err := nodeCmd.CombinedOutput()

	result := strings.TrimSpace(string(out))
	if err != nil {
		return fmt.Sprintf("[engine] Node.js execution failed!\nError: %s\nOutput:\n%s",
			err.Error(),
			string(out),
		), nil
	}
	return result, nil
}
