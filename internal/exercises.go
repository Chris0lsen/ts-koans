package internal

type Exercise struct {
	ID          string
	title       string
	description string
	StarterCode string
	TestScript  string
}

func (e Exercise) Title() string       { return e.title }
func (e Exercise) Description() string { return e.description }
func (e Exercise) FilterValue() string { return e.title }

func Exercises() []Exercise {
	var exercises = []Exercise{
		{
			title:       "Type Annotations",
			description: "Fix the function signature so that add accepts two numbers and returns a number.",
			StarterCode: `// Fix the types!
export function add(a, b): ??? {
  return a + b;
}
`,
			TestScript: `
if (typeof exports.add !== "function") throw new Error("Not a function");
if (exports.add(2, 3) !== 5) throw new Error("add(2, 3) !== 5");
`,
		},
		{
			title:       "Type Inference",
			description: "What type does TypeScript infer for the return value? Add an explicit return type.",
			StarterCode: `export function double(x: number) {
  return x * 2;
}
`,
			TestScript: `
if (exports.double(5) !== 10) throw new Error("double(5) failed");
`,
		},
		{
			title:       "Union Types",
			description: "Make 'describe' accept a string or a number and return a string.",
			StarterCode: `export function describe(input: ???): string {
  return "value: " + input;
}
`,
			TestScript: `
if (exports.describe(5) !== "value: 5") throw new Error();
if (exports.describe("hi") !== "value: hi") throw new Error();
`,
		},
		{
			title:       "Interfaces",
			description: "Define the User interface and use it as the type for the function argument.",
			StarterCode: `// Fill in the interface and function argument type
interface User {
  // ...
}

export function greet(user): string {
  return "Hello, " + user.name;
}
`,
			TestScript: `
if (exports.greet({ name: "Ada", age: 36 }) !== "Hello, Ada") throw new Error();
`,
		},
		{
			title:       "Type Aliases",
			description: "Define a type alias for a union type (yes or no) and use it as the function argument.",
			StarterCode: `// Define a type alias and use it in the function
export function accept(input: ???): boolean {
  return input === "yes";
}
`,
			TestScript: `
if (exports.accept("yes") !== true) throw new Error();
if (exports.accept("no") !== false) throw new Error();
`,
		},
		{
			title:       "Literal Types",
			description: "Make the function only accept 'start' or 'stop' as argument values.",
			StarterCode: `export function command(input: ???): string {
  return "got " + input;
}
`,
			TestScript: `
if (exports.command("start") !== "got start") throw new Error();
if (exports.command("stop") !== "got stop") throw new Error();
`,
		},
		{
			title:       "keyof Operator",
			description: "Use keyof to constrain k so only valid keys of Person are allowed.",
			StarterCode: `interface Person {
  name: string;
  age: number;
}

export function getValue(obj: Person, k: ???): string | number {
  return obj[k];
}
`,
			TestScript: `
const p = { name: "Bob", age: 30 };
if (exports.getValue(p, "name") !== "Bob") throw new Error();
if (exports.getValue(p, "age") !== 30) throw new Error();
`,
		},
		{
			title:       "typeof Operator",
			description: "Use typeof to type the function argument based on the 'template' object.",
			StarterCode: `const template = { id: 1, label: "Test" };
// Use typeof to type the argument
export function label(obj: ???): string {
  return obj.label;
}
`,
			TestScript: `
if (exports.label({ id: 2, label: "hi" }) !== "hi") throw new Error();
`,
		},
		{
			title:       "Generics",
			description: "Make the identity function generic so it works with any type.",
			StarterCode: `export function identity(x): ??? {
  return x;
}
`,
			TestScript: `
if (exports.identity(3) !== 3) throw new Error();
if (exports.identity("a") !== "a") throw new Error();
`,
		},
		{
			title:       "Optional Properties",
			description: "Add a type to 'cfg' so that timeout is optional.",
			StarterCode: `interface Config {
  url: string;
  timeout?: number;
}
// Fix the argument type
export function connect(cfg): string {
  return cfg.url + (cfg.timeout ? " with timeout" : "");
}
`,
			TestScript: `
if (exports.connect({ url: "a" }) !== "a") throw new Error();
if (exports.connect({ url: "b", timeout: 100 }) !== "b with timeout") throw new Error();
`,
		},
	}

	return exercises
}
