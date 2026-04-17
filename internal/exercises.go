package internal

import "github.com/charmbracelet/lipgloss"

type Exercise struct {
	ID             string
	title          string
	description    string
	info           string
	StarterCode    string
	TestScript     string
	Label          string
	FunctionName   string
	TypeAssertions string
}

func (e Exercise) Title() string {
	if e.Label != "" {
		return e.Label
	}
	return e.title
}
func (e Exercise) Description() string { return e.description }
func (e Exercise) Info() string        { return e.info }
func (e Exercise) FilterValue() string { return e.title }

// Styles for the info panel
var (
	kw   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	code = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	bold = lipgloss.NewStyle().Bold(true)
)

func Exercises() []Exercise {
	var exercises = []Exercise{
		{
			ID:          "primitives-string",
			title:       "Primitives: string",
			Label:       "",
			description: "An entity can have the type `string`",
			info:        `A string is a sequence of characters, or even a single character. For instance, "Linji" is a string, as is "a". Even "" is a string, albeit an empty one. Got something that looks like a number, but it's wrapped in quotes, like "123"? That's a string too!`,
			StarterCode: `const monk: ??? = "Linji";`,
			TestScript: `
if (typeof monk !== "string") throw new Error("monk should be a string, got typeof monk: " + typeof monk + ", value: " + monk);

if (monk !== "Linji") throw new Error("monk should be 'Linji'");
`,
			TypeAssertions: `
// monk should be of type string
type _Check = Assert<IsType<typeof monk, string>>;
`},
		{
			ID:          "primitives-number",
			title:       "Primitives: number",
			Label:       "",
			description: "An entity can have the type `number`",
			info:        `A number in TS (or JS) can represent both integers and floating-point values. For example, 1, -5, and 3.14 are all numbers. TS also supports special numeric values like Infinity and NaN (Not a Number). However, TS does not have separate types for integers and floats; they are all just 'number'.`,
			StarterCode: `const handsClapping: ??? = 1;`,
			TestScript: `
if (typeof handsClapping !== "number") throw new Error("handsClapping should be a number");
if (handsClapping !== 1) throw new Error("handsClapping should be 1");
`,
			TypeAssertions: `
// handsClapping should be of type number
type _Check = Assert<IsType<typeof handsClapping, number>>;
`,
		},
		{
			ID:          "primitives-boolean",
			title:       "Primitives: boolean",
			Label:       "",
			description: "An entity can have the type `boolean`",
			info:        `A boolean represents a logical entity that can be either true or false. It's commonly used in conditional statements and logical operations.`,
			StarterCode: `const nature: ??? = Boolean(false);`,
			TestScript: `
if (typeof nature !== "boolean") throw new Error("nature should be a boolean");
if (nature !== false) throw new Error("nature should be false");
`,
			TypeAssertions: `
// nature should be of type boolean
type _Check = Assert<IsType<typeof nature, boolean>>;
`},
		{
			ID:          "primitives-bigint",
			title:       "Primitives: bigint",
			Label:       "",
			description: "A very large entity can have the type `bigint`",
			info:        `A bigint represents an integer with arbitrary precision. It's useful for working with very large numbers that exceed the safe integer limit for the "number" type. Like big integers! You can create a bigint by appending 'n' to the end of an integer literal, or by using the BigInt constructor.`,
			StarterCode: `const tremendous: ??? = BigInt(100);`,
			TestScript: `
if (typeof tremendous !== "bigint") throw new Error("tremendous should be a bigint");
if (tremendous !== BigInt(100)) throw new Error("tremendous should be BigInt(100)");
`,
			TypeAssertions: `
// tremendous should be of type bigint
type _Check = Assert<IsType<typeof tremendous, bigint>>;
`},
		{
			ID:          "primitives-symbol",
			title:       "Primitives: Symbol",
			Label:       "",
			description: "A unique entity can be created with the special function `Symbol()`, and its type is `symbol`.",
			info:        `A Symbol is a unique and immutable primitive value. This means only one of a kind can exist in your program! Also, you cannot loop over the properties of a Symbol, and they are not included in JSON.stringify output. They're often used as unique keys for object properties to avoid name collisions.`,
			StarterCode: `// A Symbol is truly unique
const theOne: ???= Symbol("Linji"); 
const theOnly: ??? = Symbol("Linji");
theOne !== theOnly`,
			TestScript: `
if (typeof theOne !== "symbol") throw new Error("theOne should be a symbol");
if (typeof theOnly !== "symbol") throw new Error("theOnly should be a symbol");
if (theOne === theOnly) throw new Error("Symbols with the same description are still unique");
`,
			TypeAssertions: `
// theOne and theOnly should be of type symbol
type _Check = Assert<IsType<typeof theOne, symbol>>;
type _Check2 = Assert<IsType<typeof theOnly, symbol>>;
`},
		{
			ID:          "unique-symbol",
			title:       "Primitives: Unique Symbol",
			Label:       "",
			description: "An even more unique entity; even its type is unique",
			info:        `A unique symbol is a subtype of symbol that represents a single, specific symbol. You can create a unique symbol using the 'unique symbol' type on a const declaration. This means that the type is not just 'symbol', but a specific, unique type that can only be assigned to itself.`,
			StarterCode: `const uniqueOne: ??? symbol = Symbol("one");
const uniqueTwo: ??? symbol = Symbol("two");

// Both are symbols, so this works:
const s: symbol = uniqueOne;

// But try swapping them — the compiler won't let you:
// const nope: typeof uniqueOne = uniqueTwo;  // Error!`,
			TestScript: `
if (typeof uniqueOne !== "symbol") throw new Error("uniqueOne should be a symbol");
if (typeof uniqueTwo !== "symbol") throw new Error("uniqueTwo should be a symbol");
if (typeof s !== "symbol") throw new Error("s should be a symbol");
`,
			TypeAssertions: `
// Both are assignable to symbol, but their types are narrower
type _Check1 = Assert<IsAssignable<typeof uniqueOne, symbol>>;
type _Check2 = Assert<IsNotType<typeof uniqueOne, symbol>>;
type _Check3 = Assert<IsAssignable<typeof uniqueTwo, symbol>>;
type _Check4 = Assert<IsNotType<typeof uniqueTwo, symbol>>;
// The unique symbols are not interchangeable
type _Check5 = Assert<IsNotType<typeof uniqueOne, typeof uniqueTwo>>;
`,
		},
		{
			ID:          "primitives-any",
			title:       "Primitives: any",
			Label:       "",
			description: "An entity can have `any` type",
			info:        `The 'any' type is a powerful escape hatch that allows you to opt out of type checking for a variable. When a variable is of type 'any', it can hold values of any type, and you can perform any operation on it without TypeScript raising an error. However, using 'any' should be done with caution, as it can lead to runtime errors if not used carefully. It's often better to use more specific types or 'unknown' when you want to allow for flexibility while still maintaining some level of type safety.`,
			StarterCode: `let anything: ??? = "one"; 
anything = 2;
anything = false;`,
			TestScript: `
if (anything !== false) throw new Error("anything should be false after assignments");
`,
			TypeAssertions: `
// anything should be of type any
type _Check = Assert<IsType<typeof anything, any>>;
`},
		{
			ID:          "primitives-null",
			title:       "Primitives: null",
			Label:       "",
			description: "An entity can have `null` type",
			info:        `The 'null' type represents the intentional absence of any object value. It's often used to indicate that a variable should be empty or have no value. This will become more useful when we learn about union types a little later.`,
			StarterCode: `let nothing: ??? = null;`,
			TestScript: `
if (nothing !== null) throw new Error("nothing should be null");
`,
			TypeAssertions: `
// nothing should be of type null
type _Check = Assert<IsType<typeof nothing, null>>;
`},
		{
			ID:          "primitives-undefined",
			title:       "Primitives: undefined",
			description: "An unset entity has the type `undefined`",
			info:        `When an entity is "undefined", it means it has been declared but not assigned a value. This is different from "null", which represents the intentional absence of any object value. In JS and TS, if you declare a variable without initializing it, it will have the value 'undefined' by default. Additionally, if you try to access a property that doesn't exist on an object, it will also return 'undefined'.`,
			StarterCode: `let unset: ??? = undefined;`,
			TestScript: `
if (unset !== undefined) throw new Error("unset should be undefined");
`,
			TypeAssertions: `
// unset should be of type undefined
type _Check = Assert<IsType<typeof unset, undefined>>;
`,
		},

		{
			ID:          "arrays",
			title:       "Arrays",
			Label:       "",
			description: "An array of entities may be defined as an Array",
			info:        `An array is an ordered collection of values. In JS, arrays can hold values of any type, and even several different types at once! In TS, however, we can specify the type of values an array can hold.`,
			StarterCode: `let anArray: ???<string> = ["one", "two"]; `,
			TestScript: `
if (!Array.isArray(anArray)) throw new Error("anArray should be an array");
if (anArray.length !== 2) throw new Error("anArray should have length 2");
if (anArray[0] !== "one" || anArray[1] !== "two") throw new Error("Array elements should match");
`,
			TypeAssertions: `
// anArray should be of type Array<string>
type _Check = Assert<IsType<typeof anArray, Array<string>>>;
`},
		{
			ID:          "readonly-arrays",
			title:       "ReadonlyArrays",
			Label:       "",
			description: "A readonly array may never change",
			info:        `A ReadonlyArray is an array that cannot be modified after its creation. This means you cannot add, remove, or change elements in the array. It's useful for ensuring that data remains immutable and preventing accidental modifications.`,
			StarterCode: `let aReadonlyArray: ???<string> = ["steadfast", "unchanging"];`,
			TestScript: `
if (aReadonlyArray.length !== 2) throw new Error("aReadonlyArray should remain length 2");
if (aReadonlyArray[0] !== "steadfast" || aReadonlyArray[1] !== "unchanging") throw new Error("Elements should be unchanged");
`,
			TypeAssertions: `
// aReadonlyArray should be of type ReadonlyArray<string>
type _Check = Assert<IsType<typeof aReadonlyArray, ReadonlyArray<string>>>;
`},
		{
			ID:          "parameter-type-annotations-string",
			title:       "Parameter Type Annotations: string",
			Label:       "",
			description: "A function can accept a string",
			info:        `In addition to defining the types of variables, we can define the types that our functions will expect! This way, the TS compiler can make sure we're only passing strings to functions that expect strings, for instance.`,
			StarterCode: `function hello(name: ???) {
  return "Hello " + name;
}`,
			TestScript: `
if (hello("world") !== "Hello world") throw new Error('hello("world") should return "Hello world"');
`,
			TypeAssertions: `
// The parameter 'name' should be of type string
type _Check = Assert<IsType<Parameters<typeof hello>[0], string>>;
`},
		{
			ID:          "parameter-type-annotations-number",
			title:       "Parameter Type Annotations: number",
			Label:       "",
			description: "A function can accept a number",
			info:        `Telling this function it will always receive a number means we can perform number-like actions on it without worry!`,
			StarterCode: `function foo(bar: ???) {
  return 100 + bar;
}`,
			TestScript: `
if (foo(23) !== 123) throw new Error('foo(23) should return 123');
`,
			TypeAssertions: `
// The parameter 'bar' should be of type number
type _Check = Assert<IsType<Parameters<typeof foo>[0], number>>;
`},
		{
			ID:          "parameter-type-annotations-boolean",
			title:       "Parameter Type Annotations: boolean",
			Label:       "",
			description: "A function can accept a boolean value",
			info:        `This little toy function is purely pedantic. Technically !!value would work with any value, not just a boolean! It's nifty shorthand to convert a value into a boolean.`,
			StarterCode: `function isTrue(value: ???) {
  return !!value ? "It is true" : "It is untrue";
}`,
			TestScript: `
if (isTrue(true) !== "It is true") throw new Error('isTrue(true) should return "It is true"');
if (isTrue(false) !== "It is untrue") throw new Error('isTrue(false) should return "It is untrue"');
`,
			TypeAssertions: `
// The parameter 'value' should be of type boolean
type _Check = Assert<IsType<Parameters<typeof isTrue>[0], boolean>>;
`},
		{
			ID:          "parameter-type-annotations-any",
			title:       "Parameter Type Annotations: any",
			Label:       "",
			description: "A function can accept `any` value",
			info:        `This is another one to be careful with. Telling the compiler to expect "any" value means it can't protect us from ourselves. Also, see that "typeof" operator? We'll play with that more later too!`,
			StarterCode: `function anything(value: ???) {
  return typeof value;
}`,
			TestScript: `
if (anything("str") !== "string") throw new Error('anything("str") should return "string"');
if (anything(123) !== "number") throw new Error('anything(123) should return "number"');
if (anything(false) !== "boolean") throw new Error('anything(false) should return "boolean"');
`,
			TypeAssertions: `
// The parameter 'value' should be of type any
type _Check = Assert<IsType<Parameters<typeof anything>[0], any>>;
`},
		{
			ID:          "parameter-type-annotations-array",
			title:       "Parameter Type Annotations: Array",
			Label:       "",
			description: "A function can accept an array of values of many types",
			info:        `Of course, we can also pass arrays as arguments. An equivalent syntax is string[]. You can use whichever you prefer!`,
			StarterCode: `function theyAreTrue(values: ???<string>) {
  return values.every(value => typeof value === "string")
}`,
			TestScript: `
if (!theyAreTrue(["a", "b", "c"])) throw new Error('theyAreTrue(["a", "b", "c"]) should return true');
if (theyAreTrue(["a", 2, "c"])) throw new Error('theyAreTrue(["a", 2, "c"]) should return false');
`,
			TypeAssertions: `
// The parameter 'values' should be of type Array<string>
type _Check = Assert<IsType<Parameters<typeof theyAreTrue>[0], Array<string>>>;
`,
		},
		{
			ID:          "return-type-annotations-string",
			title:       "Return Type Annotations: string",
			Label:       "",
			description: "A function can return a string",
			info:        `Not only can we type the values going into a function - we can also define what should be returned.`,
			StarterCode: `function stringReturner(value: string): ??? {
  return value.toUpperCase()
}`,
			TestScript: `
if (stringReturner("hello") !== "HELLO") throw new Error('stringReturner("hello") should return "HELLO"');
`,
			TypeAssertions: `
// The return type should be string
type _Check = Assert<IsType<ReturnType<typeof stringReturner>, string>>;
`},
		{
			ID:          "return-type-annotations-number",
			title:       "Return Type Annotations: number",
			Label:       "",
			description: "A function can return a number",
			info:        `TypeScript is all about keeping us safe from ourselves. If we tried to return something other than a number here, the compiler would warn us.`,
			StarterCode: `function numberReturner(value: number): ??? {
  return value * 2;
}`,
			TestScript: `
if (numberReturner(21) !== 42) throw new Error('numberReturner(21) should return 42');
`,
			TypeAssertions: `
// The return type should be number
type _Check = Assert<IsType<ReturnType<typeof numberReturner>, number>>;
`},
		{
			ID:          "return-type-annotations-boolean",
			title:       "Return Type Annotations: boolean",
			Label:       "",
			description: "A function can return a boolean value",
			info:        `Remember to take breaks! Drink some water, stretch!`,
			StarterCode: `function boolReturner(value: boolean): ??? {
  return !value;
}`,
			TestScript: `
if (boolReturner(true) !== false) throw new Error('boolReturner(true) should return false');
if (boolReturner(false) !== true) throw new Error('boolReturner(false) should return true');
`,
			TypeAssertions: `
// The return type should be boolean
type _Check = Assert<IsType<ReturnType<typeof boolReturner>, boolean>>;
`},
		{
			ID:          "return-type-annotations-any",
			title:       "Return Type Annotations: any",
			Label:       "",
			description: "A function can return any value",
			info:        `Just because you can, does not mean you should.`,
			StarterCode: `function anyReturner(value: any): ??? {
  return value;
}`,
			TestScript: `
if (anyReturner(42) !== 42) throw new Error("anyReturner(42) should return 42");
if (anyReturner("foo") !== "foo") throw new Error('anyReturner("foo") should return "foo"');
`,
			TypeAssertions: `
// The return type should be any
type _Check = Assert<IsType<ReturnType<typeof anyReturner>, any>>;
`},
		{
			ID:          "return-type-annotations-void",
			title:       "Return Type Annotations: void",
			Label:       "",
			description: "A function can return to the void",
			info:        `Sometimes we value functions for thir side effects, and ask for nothing in return.`,
			StarterCode: `function voidReturner(value: any): ??? {
  return;
}`,
			TestScript: `
const result = voidReturner(123);
if (typeof result !== "undefined") throw new Error("voidReturner should return undefined");
`,
			TypeAssertions: `
// The return type should be void
type _Check = Assert<IsType<ReturnType<typeof voidReturner>, void>>;
`},
		{
			ID:          "anonymous-functions",
			title:       "Anonymous Functions",
			Label:       "",
			description: "Though nameless, anonymous functions must still abide by typing rules",
			info:        `Even though there's no specific type annotation here, the compiler sees what you're doing. Many, though, will say that Explicit is better than Implicit.`,
			StarterCode: `const monks = ["Zhaozhou", "Huineng", ???]
monks.forEach((monk) => {
  console.log(monk + " practices typescript")
})`,
			TestScript: `
if (monks.length !== 3) throw new Error("There should be three monks in the array");
if (!monks.includes("Zhaozhou")) throw new Error("Missing Zhaozhou");
if (!monks.includes("Huineng")) throw new Error("Missing Huineng");
`,
			TypeAssertions: `
// monks should be an array of strings
type _Check = Assert<IsType<typeof monks, string[]>>;
`},
		{
			ID:          "object-types",
			title:       "Object Types",
			Label:       "",
			description: "A function can accept an object of a given shape",
			info:        `Okay, technically there are a few JS quirks that could come into play here. Like "adding" a number to a string results in a concatenation operation. But let's not stray from the path.`,
			StarterCode: `function foo(value: {bar: string, baz: ???}): string {
  return value.bar + value.baz;
}`,
			TestScript: `
if (foo({ bar: "one", baz: "two" }) !== "onetwo") throw new Error('foo({ bar: "one", baz: "two" }) should return "onetwo"');
`,
			TypeAssertions: `
// The parameter 'value' should have bar: string and baz: string
type _Check = Assert<IsType<Parameters<typeof foo>[0], { bar: string; baz: string }>>;
`},
		{
			ID:          "object-types-readonly",
			title:       "Object Types: readonly",
			Label:       "",
			description: "A function can accept an object with immutable properties",
			info:        `Our friend readonly is back!`,
			StarterCode: `function foo(value: {bar: string, ??? baz: string}): void {
  value.bar = "I can change";
  value.baz !== "I cannot";
}`,
			TestScript: `
// Should be able to reassign bar, but NOT baz
let obj = { bar: "abc", baz: 42 };
foo(obj);
obj.bar = "xyz"; // allowed
let failed = false;
try {
  // @ts-expect-error
  obj.baz = 99; // should error at compile time if baz is readonly
} catch { failed = true; }
if (!("baz" in obj)) throw new Error("obj should have baz property");
`,
			TypeAssertions: `
// 'baz' should be readonly
type _Check = Assert<IsType<{ bar: string, readonly baz: string }, Parameters<typeof foo>[0]>>;
`},
		{
			ID:          "optional-properties",
			title:       "Optional Properties",
			Label:       "",
			description: "A function may accept questionable properties",
			info:        `If you attempt to access the value of an optional property, you'll get undefined.`,
			StarterCode: `// Let foo accept an optional property called bar
function foo(value: { ???: string }): boolean {
  // bar might be missing!
  if ("bar" in value) {
    return typeof value.bar === "string";
  }
  return true; // If bar is missing, that's ok
}`,
			TestScript: `
// Should allow missing bar
if (foo({}) !== true) throw new Error("foo({}) should return true");
// Should allow bar as string
if (foo({ bar: "baz" }) !== true) throw new Error('foo({ bar: "baz" }) should return true');
`,
			TypeAssertions: `
// bar should be optional. Optional parameters are declared with a ?
type _Check = Assert<IsType<{ bar?: string }, Parameters<typeof foo>[0]>>;
`},
		{
			ID:          "union-types",
			title:       "Union Types",
			Label:       "",
			description: "Several types may exist in harmony with `|`",
			info:        `The union operator | allows us to say that a value can be one of several types. It's common in TS to reach for union types instead of "any" or an enum (which we'll discuss later).`,
			StarterCode: `let something: string ??? number;
something = "Hello";
something = 100;`,
			TestScript: `
if (typeof something !== "number") throw new Error("After assignment, something should be a number");
something = "world";
if (typeof something !== "string") throw new Error("After assignment, something should be a string");
`,
			TypeAssertions: `
// something should be string | number
type _Check = Assert<IsAssignable<typeof something, string | number>>;
`},
		{
			ID:          "union-types-narrowing",
			title:       "Union Types: Narrowing",
			Label:       "",
			description: "One may narrow the union. The compiler will deduce the most specific type.",
			info:        `By checking the type of foo at runtime, we can "narrow" its type and guarantee safety within a given branch.`,
			StarterCode: `function narrow(foo: number | string): true {
  if (typeof foo === "string") {
    // In this branch, TS knows foo is a string
    return (typeof foo === "string") as true;
  } else {
   return (typeof foo === ???) as true;
  }
}`,
			TestScript: `
if (narrow("hello") !== true) throw new Error('narrow("hello") should return true');
if (narrow(42) !== true) throw new Error('narrow(42) should return true');
`,
			TypeAssertions: `
// The parameter should accept number or string
type _CheckParam = Assert<IsType<Parameters<typeof narrow>[0], string | number>>;
`},
		{
			ID:          "union-types-nullable",
			title:       "Union Types: Nullable",
			Label:       "",
			description: "A common use of union types is to represent nullable values",
			info:        `By including null in the union, we can represent values that might be absent. This is often more precise than using "any" and allows us to take advantage of TypeScript's type checking.`,
			StarterCode: `function greet(name: string | ???): string {
			  if (name === null) {
			    return "Hello, monk!";
			  } else {	
			    return "Hello, " + name + "!";
			  }
}`,
			TestScript: `
if (greet("Alice") !== "Hello, Alice!") throw new Error('greet("Alice") should return "Hello, Alice!"');
if (greet(null) !== "Hello, monk!") throw new Error('greet(null) should return "Hello, monk!"');
`,
		},
		{
			ID:          "nullish-coalescing",
			title:       "Nullish Coalescing",
			Label:       "",
			description: "The nullish coalescing operator `??` can be used to provide a default value when dealing with nullable types",
			info:        `Why use ` + code.Render("??") + ` instead of ` + code.Render("||") + `? It's a good question. The ` + code.Render("||") + ` operator will return the right-hand side if the left-hand side is falsy, which includes values like 0, "", and false. This can lead to unintended consequences if you want to allow those values. The ` + code.Render("??") + ` operator, on the other hand, only returns the right-hand side if the left-hand side is null or undefined, making it a safer choice for providing default values when dealing with nullable types.`,
			StarterCode: `function greet(name: string | null): string {
			  const actualName = name ??? "monk";
			  return "Hello, " + actualName + "!";
			}`,
			TestScript: `
if (greet("Alice") !== "Hello, Alice!") throw new Error('greet("Alice") should return "Hello, Alice!"');
if (greet(null) !== "Hello, monk!") throw new Error('greet(null) should return "Hello, monk!"');
`,
		},
		{
			ID:          "optional-chaining",
			title:       "Optional Chaining",
			Label:       "",
			description: "The optional chaining operator `?.` can be used to safely access properties on nullable types",
			info:        `If an object is null or undefined, the ` + code.Render("?.") + ` operator will short-circuit and return undefined instead of throwing an error. This is especially useful when dealing with deeply nested objects or optional properties.`,
			StarterCode: `type Monk = {
			  name: string;
			  mentor?: Monk;
}
function getMentorName(monk: Monk): string {
  return monk.mentor???name ?? "No mentor";
}`,
			TestScript: `
const xingsi = { name: "Xingsi" };
const huineng = { name: "Huineng", mentor: xingsi };
if (getMentorName(huineng) !== "Xingsi") throw new Error("Huineng's mentor should be Xingsi");
if (getMentorName(xingsi) !== "No mentor") throw new Error("Xingsi should have no mentor");
`,
		},
		{
			ID:          "as-const",
			title:       "as const",
			Label:       "",
			description: "The `as const` assertion can be used to make an object literal's properties readonly and its values literal types",
			info:        `Remember ` + bold.Render("narrowing") + `? The ` + code.Render("as const") + ` assertion is a way to tell the compiler to infer the narrowest type for an object literal. It makes all properties readonly and infers literal types for the values.`,
			StarterCode: `const monk = {
			  name: "Linji",
			  age: 800
} as ???;`,
			TestScript: `
if (monk.name !== "Linji") throw new Error("monk's name should be 'Linji'");
if (monk.age !== 800) throw new Error("monk's age should be 800");
let failed = false;
try {
  // @ts-expect-error
  monk.name = "Huineng"; // should error at compile time if name is readonly
} catch { failed = true; }
if (!("name" in monk)) throw new Error("monk should have name property");
`,
		},
		{
			ID:          "discriminated-unions",
			title:       "Discriminated Unions",
			Label:       "",
			description: "A common pattern is to use a literal property to discriminate between types in a union",
			info:        `This one isn't a specific operator - it's more of a feature of the TS compiler. By providing a property common to all types in a union, we can let the compiler narrow the type based on thta property's value!`,
			StarterCode: `type Circle {
			    kind: 'circle';
				radius: number;
			}
type Square {
			kind: 'square';
			length: number;
			}
type Shape = ???;
// Create a circle and a square, using the 'kind' property to discriminate between them
function getArea(shape: Shape) {
  switch (shape.kind) {
    case "circle":
      return Math.PI * shape.radius ** 2; // TypeScript knows shape is Circle here
    case "square":
      return shape.length ** 2; // TypeScript knows shape is Square here
  }
}

			`,
			TestScript: `const myCircle = { kind: 'circle', radius: 5 };
const mySquare = { kind: 'square', length: 5 };
if (myCircle.kind !== "circle") throw new Error("myCircle should have kind 'circle'");
if (mySquare.kind !== "square") throw new Error("mySquare should have kind 'square'");
getArea(myCircle);
getArea(mySquare);
`,
		},
		{
			ID:          "type-aliases-object-types",
			title:       "Type Aliases: Object Types",
			Label:       "",
			description: "One may define a `type` as an object",
			info:        `The power of types is that we can define custom types with any shape!`,
			StarterCode: `??? MyType = {
  foo: string;
  bar: number;
}
const val: MyType = { foo: "hi", bar: 123 };  
`,
			TestScript: `
if (val.foo !== "hi") throw new Error("foo property should be 'hi'");
if (val.bar !== 123) throw new Error("bar property should be 123");
`,
			TypeAssertions: `
// MyType should be the object type with foo: string and bar: number
type _Check = Assert<IsType<MyType, { foo: string; bar: number }>>;
`},
		{
			ID:          "type-aliases-union-types",
			title:       "Type Aliases: Union Types",
			Label:       "",
			description: "A type may be the union of other types",
			info:        `Now we're combining concepts; you can use your type aliases in unions. Suppose you want to allow both people and dogs to access your website. Your login function might accept a union of ` + code.Render("Person") + ` and ` + code.Render("Dog") + ` types!`,
			StarterCode: `type MyType = {
  foo: string;
  bar: number;
}
type MyTypeOrNumber = ???;
let myVar: MyTypeOrNumber = 100`,
			TestScript: `
myVar = { foo: "baz", bar: 123 };
if (typeof myVar === "number" && myVar !== 100) throw new Error("myVar as number should be 100");
if (typeof myVar === "object" && myVar.foo !== "baz") throw new Error("myVar as MyType should have foo='baz'");
`,
			TypeAssertions: `
// MyTypeOrNumber should be MyType | number
type _Check = Assert<IsType<MyTypeOrNumber, MyType | number>>;
`},
		{
			ID:          "type-aliases-extending",
			title:       "Type Aliases: Extending",
			Label:       "",
			description: "A type may be extended with `&`",
			info:        `The intersection operator & allows us to combine types to create new ones. This is often used to extend an existing type with new properties. For instance, if we have a ` + code.Render("Person") + ` type, we can create a ` + code.Render("Monk") + ` type that includes all the properties of ` + code.Render("Person") + ` and adds some new ones!`,
			StarterCode: `type Person = {
  name: string;
}
type Monk = Person ??? {
  isMeditating: boolean;
}
const m: Monk = { name: "Linji", isMeditating: true };  
`,
			TestScript: `
if (m.name !== "Linji") throw new Error("Monk should have correct name");
if (!m.isMeditating) throw new Error("Monk should have isMeditating property true");
`,
			TypeAssertions: `
// Monk should be the intersection of Person and { isMeditating: boolean }
type _Check = Assert<IsType<Monk, Person & { isMeditating: boolean }>>;
`},
		{
			ID:          "type-aliases-immutability",
			title:       "Type Aliases: Immutability",
			Label:       "",
			description: "A type may not change after its creation",
			info:        `Once a type alias is declared, it cannot be redeclared or changed (but it can be extended). If you're in a position where you feel like you need to change a type, you might want to be using interfaces instead - or maybe you need to re-think your model!`,
			StarterCode: `// The below code will not compile.
type Constancy = {
  foo: boolean
}

// Change the name of the second type to make the code compile.
type Constancy = {
  bar: boolean
}`,
			TestScript: `
// This koan is about the fact that type aliases cannot be redeclared.
// No runtime test needed.
`,
			TypeAssertions: `
// Only one type Constancy should exist
type _Check1 = Assert<IsType<Constancy, { foo: boolean }>>;
`},
		{
			ID:          "interfaces",
			title:       "Interfaces",
			Label:       "",
			description: "An interface is very similar to a type",
			info:        `Interfaces are extremely similar to type aliases. In fact, for object types, they are almost interchangeable. Interfaces can be altered after declaration, while types cannot. It's conventional to use interfaces for most object types, and to use type aliases for things like unions and intersections, but you may walk your own path!`,
			StarterCode: `??? MyInterface {
  foo: string;
  bar: number;
}
const obj: MyInterface = { foo: "hello", bar: 123 };
`,
			TestScript: `
if (obj.foo !== "hello") throw new Error("foo should be 'hello'");
if (obj.bar !== 123) throw new Error("bar should be 123");
`,
			TypeAssertions: `
// MyInterface should match the object type { foo: string; bar: number }
type _Check = Assert<IsType<MyInterface, { foo: string; bar: number }>>;
`},
		{
			ID:          "interfaces-extending",
			title:       "Interfaces: Extending",
			Label:       "",
			description: "An interface can be extended as well, with `extends`",
			info:        `Just like type aliases, interfaces can also be extended to create new interfaces. This is done using the ` + code.Render("extends") + ` keyword. When an interface extends another, it inherits all of its properties and can also add new ones. This is a common way to create more specific types based on more general ones.`,
			StarterCode: `interface Person {
  name: string;
}
interface Monk ??? Person {
  isMeditating: boolean
}
const m: Monk = { name: "Huineng", isMeditating: true };
`,
			TestScript: `
if (m.name !== "Huineng") throw new Error("Monk should have correct name");
if (!m.isMeditating) throw new Error("Monk should have isMeditating property true");
`,
			TypeAssertions: `
// Monk should extend Person and add isMeditating: boolean
type _Check = Assert<IsType<Monk, {name: string, isMeditating: boolean }>>;
`},
		{
			ID:          "interfaces-redefining",
			title:       "Interfaces: Redefining",
			Label:       "",
			description: "An interface can be redefined freely, merging the declarations",
			info:        `This is a powerful and confusing feature of interfaces. If you attempt to redeclare an interface, TS will instead merge together all existing declarations of that interface.`,
			StarterCode: `interface MyInterface {
  foo: string;
}
??? MyInterface {
  bar: number;
}
const obj: MyInterface = { foo: "hi", bar: 5 };
`,
			TestScript: `
if (obj.foo !== "hi") throw new Error("foo should be 'hi'");
if (obj.bar !== 5) throw new Error("bar should be 5");
`,
			TypeAssertions: `
// MyInterface should include both foo and bar
type _Check = Assert<IsType<MyInterface, { foo: string; bar: number }>>;
`},
		{
			ID:          "tuples",
			title:       "Tuples",
			Label:       "",
			description: "A tuple is an array that knows its shape and size",
			info:        `A tuple is a special type of array, of fixed length and order, where each element is explicitly typed.`,
			StarterCode: `function foo(myTuple: [string, ???]): true {
  return (typeof myTuple[0] === "string"
  && typeof myTuple[1] === "number") as true;
}`,
			TestScript: `
if (!foo(["a", 1])) throw new Error('foo(["a", 1]) should return true');
`,
			TypeAssertions: `
// myTuple should be a tuple [string, number]
type _Check = Assert<IsType<Parameters<typeof foo>[0], [string, number]>>;
`},
		{
			ID:          "readonly-tuples",
			title:       "Readonly Tuples",
			Label:       "",
			description: "A tuple can be readonly",
			info:        `A tuple can be readonly. You might need this one day.`,
			StarterCode: `function foo(myTuple: ??? [string, number]): void {
  console.log(myTuple[0] + " will always be a string")
}
const tuple: Readonly<[string, number]> = ["foo", 42] as const;
`,
			TestScript: `
// Should accept a readonly tuple
foo(tuple);
// Should error if trying to mutate (compile-time error if parameter is readonly)
let failed = false;
try {
  // @ts-expect-error
  tuple[0] = "bar";
} catch { failed = true; }
if (!("length" in tuple)) throw new Error("Tuple should still exist");
`,
			TypeAssertions: `
// myTuple should be a readonly tuple [string, number]
type _Check = Assert<IsType<Parameters<typeof foo>[0], Readonly<[string, number]>>>;
`},
		{
			ID:          "promises",
			title:       "Promises",
			Label:       "",
			description: "There exists a special `Promise` type for functions that return promises",
			info:        `Asynchronous JS is so common that TS has a built-in type for it. By providing a type to the Promise utility type, we can tell the compiler what the promise will resolve to.`,
			StarterCode: `async function foo(): ???<number> {
  return 100;
}`,
			TestScript: `
foo().then(val => {
  if (val !== 100) throw new Error("Promise should resolve to 100");
});
`,
			TypeAssertions: `
// The return type should be Promise<number>
type _Check = Assert<IsType<ReturnType<typeof foo>, Promise<number>>>;
`},
		{
			ID:          "type-assertions-as",
			title:       "Type Assertions: as",
			Label:       "",
			description: "Sometimes you may need to tell the compiler what type to expect",
			info:        `You are a human. There might be a time when you know something your computer doesn't. On these days, you can instruct the compiler to expect a certain type.`,
			StarterCode: `type SometimesANumber = number | string
function numberReturner(flag: boolean): SometimesANumber {
  return flag ? "Hello" : 100;
}
const myNum: number = numberReturner(false) ??? number;`,
			TestScript: `
if (myNum !== 100) throw new Error("myNum should be 100");
`,
			TypeAssertions: `
// myNum should be a number, asserted from SometimesANumber
type _Check = Assert<IsType<typeof myNum, number>>;
`},
		{
			ID:          "literal-types",
			title:       "Literal Types",
			Label:       "",
			description: "A type can be literally `anything`",
			info:        `A literal type is a type that represents a specific value. In this case, the variable ` + code.Render("anything") + ` can only have the value ` + code.Render("anything") + `. This might be useful if you have, say, a union of string literals and you want to ensure a variable is one of those specific strings.`,
			StarterCode: `let anything: "anything" = ???`,
			TestScript: `
if (anything !== "anything") throw new Error('anything should be "anything"');
`,
			TypeAssertions: `
// anything should only allow the value "anything"
type _Check = Assert<IsType<typeof anything, "anything">>;
`},
		{
			ID:          "literal-types-unions-of-strings",
			title:       "Literal Types: Unions Of strings",
			Label:       "",
			description: "A type can be a union of strings",
			info:        `Hey, we just talked about this! Maybe you want a variable to only accept one of a few possible values. A union of string literals is a way to do that.`,
			StarterCode: `type ManyThings = "one" | "another" | ???
let thing = "a secret third thing"`,
			TestScript: `
thing = "another";
thing = "a secret third thing";
if (thing !== "a secret third thing") throw new Error('thing should be "a secret third thing"');
`,
			TypeAssertions: `
// ManyThings should include the string "a secret third thing"
type _Check = Assert<IsAssignable<"a secret third thing", ManyThings>>;
`},

		{
			ID:          "literal-types-unions-of-numbers",
			title:       "Literal Types: Unions Of numbers",
			Label:       "",
			description: "A type can be a union of numbers",
			info:        `As with strings, we can create unions of number literals. I think you probably see where this is headed.`,
			StarterCode: `type ManyNumbers = 1 | 2 | ???
let myNumber = 100`,
			TestScript: `
let n = 1;
n = 2;
n = 100;
if (n !== 100) throw new Error("n should be 100");
`,
			TypeAssertions: `
// ManyNumbers should include the number 100
type _Check = Assert<IsAssignable<100, ManyNumbers>>;
`,
		},
		{
			ID:          "literal-types-as-literal",
			title:       "Literal Types: as Literal",
			Label:       "",
			description: "Literal types may require assertion",
			info:        `Sometimes TS won't be able to infer that a variable with a literal union type is actually a specific type. You can be assertive.`,
			StarterCode: `function foo(value: "bar" | "baz"): void {}
const myValue = "bar"
// Coerce the compiler with a single operator
foo(myValue ??? "bar")`,
			TypeAssertions: `
// Should be assignable to "bar". Use "as bar" to make the compiler pass.
type _Assert = Assert<IsType<typeof myValue, "bar">>;
`,
		},
		{
			ID:          "enums-number",
			title:       "Enums: number",
			Label:       "",
			description: "Enums are sets of named constants that auto-increment",
			info:        `Enums are a way to define a set of named constants. By default, they auto-increment from 0, but you can also assign specific values. Here's a funny TS quirk: most TS types don't actually generate any JS code - they're just for the compiler. Enums, on the other hand, do generate real JS objects, which is why they have some unique behaviors.`,
			StarterCode: `enum Colors {
  Red = 0,
  Green,
  Blue,
}
Colors.Blue === ???`,
			TestScript: `
if (Colors.Blue !== 2) throw new Error("Colors.Blue should be 2");
`,
			TypeAssertions: `
// Colors.Blue should have Colors.Blue as its type
type _Check = Assert<IsType<typeof Colors.Blue, Colors.Blue>>;
`},
		{
			ID:          "enums-string",
			title:       "Enums: string",
			Label:       "",
			description: "Enums can have string values",
			info:        `Enums can also have string values. Unlike number enums, string enums do not auto-increment. That would be unreasonable.`,
			StarterCode: `enum Colors {
  Red = "RED",
  Green = "GREEN",
  Blue = "BLUE",
}
Colors.Blue === ???`,
			TestScript: `
if (Colors.Blue !== "BLUE") throw new Error('Colors.Blue should be "BLUE"');
`,
			TypeAssertions: `
// Colors.Blue should have Colors.Blue as its type
type _Check = Assert<IsType<typeof Colors.Blue, Colors.Blue>>;
`},
		{
			ID:          "type-guards-typeof",
			title:       "Type Guards: typeof",
			Label:       "",
			description: "`typeof` can be used in expressions or in types",
			info:        `The ` + code.Render("typeof") + ` operator is a powerful tool that we've seen throughout these exercises. It can be used in expressions to check the type of a variable at runtime, and it can also be used in type assertions to infer types based on the value of a variable.`,
			StarterCode: `let foo = "foo";
let bar: ??? foo;
bar = "bar"
typeof foo === typeof bar;`,
			TestScript: `
if (typeof bar !== "string") throw new Error("bar should be a string");
if (bar !== "foo") bar = "bar"; // allowed
`,
			TypeAssertions: `
// bar should be of type string, inferred from typeof foo
type _Check = Assert<IsType<typeof bar, typeof foo>>;
`},
		{
			ID:          "narrowing-in",
			title:       "Narrowing: in",
			Label:       "",
			description: "`in` can be used to narrow types",
			info:        `The ` + code.Render("in") + ` operator can be used to check if a property exists in an object. This is useful for narrowing types when you have a union of object types.`,
			StarterCode: `type PersonType = {
  name: string
}
type ObjectType = {
  foo: string
}
function typeDecider(thing: PersonType | ObjectType): string {
  return "name" ??? thing ? thing.name : thing.foo;
}`,
			TestScript: `
if (typeDecider({ name: "Linji" }) !== "Linji") throw new Error('typeDecider({ name: "Linji" }) should return "Linji"');
if (typeDecider({ foo: "bar" }) !== "bar") throw new Error('typeDecider({ foo: "bar" }) should return "bar"');
`,
			TypeAssertions: `
// typeDecider should return "Person" or "Object" as a string
type _Check = Assert<IsType<ReturnType<typeof typeDecider>, string>>;
`},
		{
			ID:          "type-predicates-is",
			title:       "Type Predicates: is",
			Label:       "",
			description: "A type predicate will tell the compiler about the type of a variable",
			info:        `Remember that sometimes you will know more than the compiler. You may use the ` + code.Render("is") + ` operator to create what is called a type predicate. It takes the form ` + code.Render("myParameterName is someType") + ` and tells the compiler that, ` + bold.Render("if") + ` the function returns true, then the parameter is of the specified type.`,
			StarterCode: `type Monk = "Linji" | "Zhaozhou"
function isPerson(value: unknown): value ??? Monk {
  return value === "Linji" || value === "Zhaozhou"
}`,
			TestScript: `
if (!isPerson("Linji")) throw new Error("isPerson('Linji') should be true");
if (!isPerson("Zhaozhou")) throw new Error("isPerson('Zhaozhou') should be true");
if (isPerson("Chris")) throw new Error("isPerson('Chris') should be false");
if (isPerson(100)) throw new Error("isPerson(100) should be false");
`,
			TypeAssertions: `
// TypeScript should know that "Linji" is a Monk after the check:
function testNarrowing(x: unknown) {
	if (isPerson(x)) {
		type _Assert = Assert<IsType<typeof x, Monk>>;
	}
}
`,
		},
		{
			ID:          "narrowing-never",
			title:       "Narrowing: never",
			Label:       "",
			description: "The `never` type represents values that never occur.",
			info:        `This is uncommon, but not rare. Some paths are forbidden.`,
			StarterCode: `function fail(message: string): ??? {
  throw new Error(message);
}`,
			TestScript: `
let threw = false;
try { fail("boom"); } catch { threw = true; }
if (!threw) throw new Error("fail() should throw");
`,
			TypeAssertions: `
// This should compile only if fail() returns never:
type _Assert = Assert<IsType<ReturnType<typeof fail>, never>>;
`,
		},
		{
			ID:          "index-signatures",
			title:       "Index Signatures",
			Label:       "",
			description: "Index signatures let you type objects with unknown `key`s, but known value types.",
			info:        `Sometimes you'll want to create object types, but you won't know the key names at compile time. Don't worry! Somebody has thought of this already.`,
			StarterCode: `interface PersonAgeMap {
    [???: string]: number
}

const ages: PersonAgeMap = {};
ages["Chris"] = 36;
ages["Linji"] = 1159;`,
			TestScript: `
if (ages["Chris"] !== 36) throw new Error("ages['Chris'] should be 36");
if (ages["Linji"] !== 1159) throw new Error("ages['Linji'] should be 1159");
`,
			TypeAssertions: `
// The following line should NOT type-check if value is not a number:
type _Check = Assert<IsType<PersonAgeMap["whatever"], number>>;
// The following line should NOT type-check if key is not a string:
type _KeyCheck = Assert<IsNotType<keyof PersonAgeMap, number>>;
`,
		},
		{
			ID:          "index-signature-unknown",
			title:       "Index Signatures & `unknown`",
			Label:       "",
			description: "The `unknown` type is a safer alternative to any.",
			info:        `Why not just use ` + code.Render("any") + `? The ` + code.Render("any") + ` type is a way to opt-out of type checking altogether. The ` + code.Render("unknown") + ` type, on the other hand, forces you to perform some kind of type check before you can use the value, making it a safer choice when you don't know the exact type of the values in your object.`,
			StarterCode: `interface AnyData {
  [key: string]: ???;
}

const record: AnyData = {};
record["foo"] = 123;
record["bar"] = "hello";`,
			TestScript: `
if (record["foo"] !== 123) throw new Error("record['foo'] should be 123");
if (record["bar"] !== "hello") throw new Error("record['bar'] should be 'hello'");
`,
			TypeAssertions: `
// The following should fail without a type assertion or guard:
type _ShouldError = Assert<IsNotType<number, typeof record["foo"]>>;

// Should allow type guard:
function useFoo(map: AnyData) {
    const val = map["foo"];
    if (typeof val === "number") {
        type _Num = Assert<IsType<typeof val, number>>;
    }
}
`,
		},
		{
			ID:          "intersection-types",
			title:       "Intersection Types",
			Label:       "",
			description: `Intersection types (using &) combine multiple types into one.`,
			info:        `We've seen this operator before - we used it to extend types. But did you know it has another use? It can be a little confusing if you're thinking about it in terms of set theory - but an intersection type in TS represents a subset of values that satisfy all of the combined types. For example, if we have a type that represents objects with a name property, and another type that represents objects with an age property, we can create an intersection type that represents objects that have both a name and an age.`,
			StarterCode: `type HasName = { name: string };
type HasAge = { age: number };

type Person = HasName ??? HasAge;

const user: Person = { name: "Hakuin", age: 256 };`,
			TestScript: `
if (user.name !== "Hakuin") throw new Error("user.name should be 'Hakuin'");
if (user.age !== 256) throw new Error("user.age should be 256");
`,
			TypeAssertions: `
// Should type-check if user is both HasName and HasAge
type _AssertName = Assert<IsAssignable<typeof user, HasName>>;
type _AssertAge = Assert<IsAssignable<typeof user, HasAge>>;
`,
		},
		{
			ID:          "generics-type-alias",
			title:       "Generics: Type Alias",
			Label:       "",
			description: `You can create reusable types with generics.`,
			info:        `"Why would I need this?" I hear you asking yourself. But it is more common than you might expect. This Box can hold anything. You might want to give it other box-like properties as well. You can do this without creating a separate type for every possible value.`,
			StarterCode: `type Box<T> = {
    value: ???
}

const numBox: Box<number> = { value: 123 }
const strBox: Box<string> = { value: "hi" }`,
			TestScript: `
if (numBox.value !== 123) throw new Error("numBox.value should be 123");
if (strBox.value !== "hi") throw new Error("strBox.value should be 'hi'");
`,
			TypeAssertions: `
// Should type-check for both number and string
type _Num = Assert<IsType<typeof numBox, Box<number>>>;
type _Str = Assert<IsType<typeof strBox, Box<string>>>;
`,
		},
		{
			ID:          "generics-function",
			title:       "Generics: Function",
			Label:       "",
			description: `Functions can also be generic!`,
			info:        `Mayhap you'll need a function that can accept and return any type, so long as they're the same type.`,
			StarterCode: `function identity<T>(value: T): ??? {
	return value;
}
	`,
			TestScript: `
if (identity(123) !== 123) throw new Error("identity(123) should return 123");
if (identity("hello") !== "hello") throw new Error('identity("hello") should return "hello"');
`,
			TypeAssertions: `
// The return type should be the same as the parameter type
type _CheckNum = Assert<IsType<ReturnType<typeof identity>, Parameters<typeof identity>[0]>>;
`,
		},
		{
			ID:          "generics-constraints",
			title:       "Generics: Constraints",
			Label:       "",
			description: `You can constrain generic types to ensure they have certain properties.`,
			info:        `The syntax can be overwhelming here. We have a generic function that takes an object of type ` + code.Render("T") + ` and a key of type ` + code.Render("K") + `. The ` + code.Render("K extends keyof T") + ` part is a constraint that says "K must be a key of T". This means that when you call ` + code.Render("getProperty") + `, the compiler will ensure that the key you provide is actually a valid key for the object you're passing in. This allows us to safely access properties on the object without risking a runtime error.`,
			StarterCode: `function getProperty<T, K extends keyof T>(obj: T, key: K): ??? {
	return obj[key];
}`,
			TestScript: `
const person = { name: "Dogen", age: 900 };
if (getProperty(person, "name") !== "Dogen") throw new Error('getProperty(person, "name") should return "Dogen"');
if (getProperty(person, "age") !== 900) throw new Error('getProperty(person, "age") should return 900');
`,
			TypeAssertions: `
// getProperty should return the correct type for a given key
const _name = getProperty({ name: "Dogen", age: 900 }, "name");
type _Check = Assert<IsType<typeof _name, string>>;
const _age = getProperty({ name: "Dogen", age: 900 }, "age");
type _Check2 = Assert<IsType<typeof _age, number>>;
`,
		},
		{
			ID:          "generics-defaults",
			title:       "Generics: Defaults",
			Label:       "",
			description: `Generic type parameters can have defaults, making them optional when using the generic.`,
			info:        `If no type argument is provided, the default type will be used. That's how defaults work! You knew that. Anyway, here's how you do it in TS. It also works for interfaces, and with multiple type parameters.`,
			StarterCode: `type Box<T = ???> = {
	value: T;
}

const defaultBox: Box = { value: "hello" };
`,
			TestScript: `
const defaultBox = { value: "hello" };
const numberBox = { value: 123 };
if (defaultBox.value !== "hello") throw new Error("defaultBox.value should be 'hello'");
if (numberBox.value !== 123) throw new Error("numberBox.value should be 123");
`,
			TypeAssertions: `
// defaultBox should be Box<string>
type _CheckDefault = Assert<IsType<typeof defaultBox, Box<string>>>;
// numberBox should still be Box<number>
type _CheckNumber = Assert<IsType<typeof numberBox, Box<number>>>;
`,
		},
		{
			ID:          "keyof",
			title:       "The keyof Keyword",
			Label:       "",
			description: `keyof returns a union of the keys of the given type`,
			info:        `This comes in handy, believe it or not. You might need to create a type that represents the keys of another type. You can combine this with generics in order to work with the keys of types you might not know at compile time! Doesn't that sound fun?`,
			StarterCode: `type User = {
    name: string;
    age: number;
    email: string;
}

type UserKeys = ???

const k1: UserKeys = "name"
const k2: UserKeys = "age"
const k3: UserKeys = "email"`,
			TestScript: `
if (k1 !== "name") throw new Error("k1 should be 'name'");
if (k2 !== "age") throw new Error("k2 should be 'age'");
if (k3 !== "email") throw new Error("k3 should be 'email'");
`,
			TypeAssertions: `
// Should only allow these keys:
type _Assert = Assert<IsType<UserKeys, "name" | "age" | "email">>;
`,
		},
		{
			ID:          "mapped-types",
			title:       "Mapped Types",
			Label:       "",
			description: "A mapped type lets you create a new type by transforming all properties of another type.",
			info:        `Just as you can map over arrays for create new arrays, you can map over types to create new types.`,
			StarterCode: `type User = {
    id: number;
    username: string;
    email: string;
}

// TODO: Make BooleanFlags so that every property of User is a boolean
type BooleanFlags = {
    [K in keyof User]: ???
}

const flags: BooleanFlags = {
    id: true,
    username: false,
    email: true
}`,
			TestScript: `
if (flags.id !== true) throw new Error("flags.id should be true");
if (flags.username !== false) throw new Error("flags.username should be false");
`,
			TypeAssertions: `
// All fields of User should be mapped to a boolean
type _Check = Assert<IsType<BooleanFlags, { id: boolean; username: boolean; email: boolean }>>;
`,
		},
		{
			ID:          "mapped-type-remove-optional",
			title:       "Mapped Type: Remove Optional Modifier",
			Label:       "",
			description: "Use a mapped type and the `-?` operator to make all properties of `MaybeUser` required.",
			info:        `There exists syntactic sugar for removing optional modifiers from properties in a mapped type.`,
			StarterCode: `type MaybeUser = {
    id?: number;
    username?: string;
    email?: string;
}

type RequiredUser = {
    [K in keyof MaybeUser]???: MaybeUser[K]
}

const u: RequiredUser = {
    id: 1,
    username: "ada",
    email: "ada@example.com"
}`,
			TestScript: `
if (!u.id || !u.username || !u.email) throw new Error("All properties should be required!");
`,
			TypeAssertions: `
// All props should be required. Use -? to remove optional properties from a mapped type
type _Check = Assert<IsType<RequiredUser, { id: number; username: string; email: string }>>;
`,
		},
		{
			ID:          "mapped-type-remove-readonly",
			title:       "Mapped Type: Remove readonly",
			Label:       "",
			description: "Use a mapped type and the `-readonly` operator to create a type where all properties are writable.",
			info:        `Just as you can subtract optional modifiers, you can subtract the readonly modifier from properties in a mapped type. Maybe you need a copy of a user that can be edited.`,
			StarterCode: `type ReadonlyUser = {
    readonly id: number;
    readonly username: string;
    readonly email: string;
}

type WritableUser = {
    ??? [K in keyof ReadonlyUser]: ReadonlyUser[K]
}

let user: WritableUser = {
    id: 1,
    username: "Ummon",
    email: "Ummon@bluecliff.com"
}

user.id = 100;`,
			TestScript: `
user.id = 2;
if (user.id !== 2) throw new Error("user.id should be mutable!");
`,
			TypeAssertions: `
// Should not be readonly; use -readonly to remove the readonly property
type _Check = Assert<IsNotReadonly<WritableUser, "id">>;
`,
		},
		{
			ID:          "utility-types-partial",
			title:       "Utility Types: `Partial`",
			Label:       "",
			info:        `There's no ` + code.Render("+?") + ` operator to make all properties optional in a mapped type, but there is a built-in utility type that does exactly that.`,
			description: "The `Partial<T>` utility type makes all properties in T optional.",
			StarterCode: `type User = {
    id: number;
    username: string;
    email: string;
}

type MaybeUser = ???

const u: MaybeUser = {};
u.id = 1;
u.username = "Bodhidharma";`,
			TestScript: `
u.email = "Bodhidharma@east.com";
if (u.id !== 1) throw new Error("id should be 1");
if (u.username !== "Bodhidharma") throw new Error("username should be 'Bodhidharma'");
if (u.email !== "Bodhidharma@east.com") throw new Error("email should be 'Bodhidharma@east.com'");
`,
			TypeAssertions: `
// Should allow all properties to be omitted
type _Check = Assert<IsType<MaybeUser, { id?: number; username?: string; email?: string }>>;
`,
		},
		{
			ID:          "utility-types-required",
			title:       "Utility Types: `Required`",
			Label:       "",
			description: "The `Required<T>` utility type makes all properties in T required (not optional).",
			info:        `This can be thought of as shorthand for using a mapped type to remove optional modifiers from all properties. It's the opposite of Partial.`,
			StarterCode: `type User = {
    id?: number;
    username?: string;
    email?: string;
}

type FullUser = ???

const u: FullUser = {
    id: 100,
    username: "Yun-men",
    email: "yun-men@sumeru.com"
}`,
			TestScript: `
if (u.id !== 100) throw new Error("id should be 100");
if (u.username !== "Yun-men") throw new Error("username should be 'Yun-men'");
if (u.email !== "yun-men@sumeru.com") throw new Error("email should be 'yun-men@sumeru.com'");
`,
			TypeAssertions: `
// Should require all properties
type _Check = Assert<IsType<FullUser, { id: number; username: string; email: string }>>;
`,
		},
		{
			ID:          "utility-pick",
			title:       "Utility Types: `Pick`",
			Label:       "",
			description: "The `Pick<T, K>` utility type creates a new type by selecting a subset of properties from T.",
			info:        `This is useful when you want to create a type that only includes a few properties from another type.`,
			StarterCode: `type User = {
    id: number;
    username: string;
    email: string;
}

// TODO: Make UserPreview with only id and username using Pick
type UserPreview = ???<User, "id" | "username">

const preview: UserPreview = {
    id: 100,
    username: "Ikkyu"
    // email: "should not exist" // should error!
}`,
			TestScript: `
if (preview.id !== 100) throw new Error("id should be 100");
if (preview.username !== "Ikkyu") throw new Error("username should be 'Ikkyu'");
`,
			TypeAssertions: `
// Should only have id and username, not email
type _Check = Assert<IsNotType<keyof UserPreview, "email">>;
type _Check2 = Assert<IsType<UserPreview, { id: number; username: string }>>;
`,
		},

		{
			ID:          "utility-types-omit",
			title:       "Utility Types: `Omit`",
			Label:       "",
			description: "Sometimes you must omit, to create something new",
			info:        `The ` + code.Render("Omit<T, K>") + ` utility type creates a new type by omitting a subset of properties from T. It's the opposite of Pick.`,
			StarterCode: `type User = {
    id: number;
    username: string;
    email: string;
    password: string;
}

type PublicUser = ???<User, "password">

const user: PublicUser = {
    id: 1,
    username: "Dongshan",
    email: "Dongshan@shouchu.com"
}`,
			TestScript: `
if (user.username !== "Dongshan") throw new Error("username should be 'Dongshan'");
if (user.email !== "Dongshan@shouchu.com") throw new Error("email should be 'Dongshan@shouchu.com'");
`,
			TypeAssertions: `
// Should not allow password
type _Check = Assert<IsNotType<keyof PublicUser, "password">>;
type _Check2 = Assert<IsType<keyof PublicUser, "id" | "username" | "email">>;
`,
		},
		{
			ID:          "utility-types-readonly",
			title:       "Utility Types: `Readonly`",
			Label:       "",
			description: "The `Readonly<T>` utility type makes all properties in T readonly.",
			info:        `This is the equivalent of using a mapped type to add the readonly modifier to all properties. It's a quick way to make an entire type immutable.`,
			StarterCode: `type User = {
	id: number;
	username: string;
	email: string;
}

const user: ???<User> = {
	id: 1,
	username: "Dongshan",
	email: "Dongshan@shouchu.com"
};
`, TestScript: `
if (user.username !== "Dongshan") throw new Error("username should be 'Dongshan'");
if (user.email !== "Dongshan@shouchu.com") throw new Error("email should be 'Dongshan@shouchu.com'");
`,
			TypeAssertions: `
// All properties should be readonly
type _Check = Assert<IsType<typeof user, Readonly<User>>>;
`,
		},
		{
			ID:          "utility-types-record",
			title:       "Utility Types: `Record`",
			Label:       "",
			description: "The `Record<K, T>` utility type constructs an object type whose keys are K and values are T.",
			info:        `Here's one you'll wind up using a lot: Records. Imagine you're waiting for an API response and know that the keys will be a specific set of strings, but you don't know how many there will be or what the values will look like. You can use Record to type this response!`,
			StarterCode: `type Page = "home" | "about" | "contact";

const pageViews: ???<Page, number> = {
	home: 1000,
	about: 500,
	contact: 200
};`,
			TestScript: `
if (pageViews.home !== 1000) throw new Error("home page should have 1000 views");
if (pageViews.about !== 500) throw new Error("about page should have 500 views");
if (pageViews.contact !== 200) throw new Error("contact page should have 200 views");
`,
			TypeAssertions: `
// Should only allow keys of Page and values of number
type _Check = Assert<IsType<typeof pageViews, Record<Page, number>>>;
`,
		},
		{
			ID:          "utility-types-returntype",
			title:       "Utility Types: `ReturnType`",
			Label:       "",
			description: "The `ReturnType<T>` utility type constructs a type consisting of the return type of function T.",
			info:        `I'll be honest, I haven't had a need for this one. But it seems cool! You can extract the return type of a function and use it elsewhere. Neat!`,
			StarterCode: `function getUser() {
	return {
		id: 1,
		username: "Shitou",
		email: "shitou@example.com"
	};
}

type User = ???<typeof getUser>
const user: User = getUser();
`, TestScript: `
if (user.username !== "Shitou") throw new Error("username should be 'Shitou'");
if (user.email !== "shitou@example.com") throw new Error("email should be 'shitou@example.com'");
`,
			TypeAssertions: `
// User should be the return type of getUser
type _Check = Assert<IsType<User, ReturnType<typeof getUser>>>;
`,
		},
		{
			ID:          "utility-types-exclude",
			title:       "Utility Types: `Exclude`",
			Label:       "",
			description: "The `Exclude<T, U>` utility type constructs a type by excluding from T all union members that are assignable to U.",
			info:        `This is like using ` + code.Render("Omit") + ` on a union type. It allows you to create a new type by excluding certain members from an existing union type.`,
			StarterCode: `type someType = string | number | boolean;
type Excluded = ???<someType, string | boolean>;

const value: Excluded = 123;`,
			TestScript: `
if (value !== 123) throw new Error("value should be 123");
`,
			TypeAssertions: `
// Should exclude string and boolean from T, leaving only number
type _Check = Assert<IsType<Excluded, number>>;
`,
		},
		{
			ID:          "utility-types-extract",
			title:       "Utility Types: `Extract`",
			Label:       "",
			description: "The `Extract<T, U>` utility type constructs a type by extracting from T all union members that are assignable to U.",
			info:        `This is the opposite of ` + code.Render("Exclude") + `. It allows you to create a new type by extracting ` + bold.Render("only the members from an existing union type") + ` that are assignable to another type. That's a verbose definition, but language is an imperfect medium.`,
			StarterCode: `type T = string | number | boolean;
type Extracted = ???<T, string | boolean>;

const value: Extracted = "hello";`,
			TestScript: `
if (value !== "hello") throw new Error("value should be 'hello'");
`,
			TypeAssertions: `
// Should extract only string and boolean from T
type _Check = Assert<IsType<Extracted, string | boolean>>;
`,
		},
	}
	return exercises
}
