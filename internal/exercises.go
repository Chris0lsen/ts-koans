package internal

type Exercise struct {
	ID             string
	title          string
	description    string
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
func (e Exercise) FilterValue() string { return e.title }

func Exercises() []Exercise {
	var exercises = []Exercise{
		{
			ID:          "primitives-string",
			title:       "Primitives: string",
			Label:       "",
			description: "An entity can have the type `string`",
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
			ID:          "primitives-any",
			title:       "Primitives: any",
			Label:       "",
			description: "An entity can have `any` type",
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
			ID:          "arrays",
			title:       "Arrays",
			Label:       "",
			description: "An array of entities may be defined as an Array",
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
			description: "Several types may exist in harmony",
			StarterCode: `let var: ???
var = "Hello";
var = 100;`,
			TestScript: `
if (typeof var !== "number") throw new Error("After assignment, var should be a number");
var = "world";
if (typeof var !== "string") throw new Error("After assignment, var should be a string");
`,
			TypeAssertions: `
// var should be string | number
type _Check = Assert<IsType<typeof var, string | number>>;
`},
		{
			ID:          "union-types-narrowing",
			title:       "Union Types: Narrowing",
			Label:       "",
			description: "One may narrow the union. The compiler will deduce the most specific type.",
			StarterCode: `function narrow(foo: number | string): true as const {
  if (typeof foo === "string") {
    return typeof foo === "string";
  } else {
   return typeof foo === ???;
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
			ID:          "type-aliases-object-types",
			title:       "Type Aliases: Object Types",
			Label:       "",
			description: "One may define a type as an object",
			StarterCode: `??? MyType = {
  foo: string;
  bar: number;
}`,
			TestScript: `
const val: MyType = { foo: "hi", bar: 123 };
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
			StarterCode: `type MyType = {
  foo: string;
  bar: number;
}
type MyTypeOrNumber = ???
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
			StarterCode: `type Person = {
  name: string;
}
type Monk = Person ??? {
  isMeditating: boolean;
}`,
			TestScript: `
const m: Monk = { name: "Linji", isMeditating: true };
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
			StarterCode: `// The below code will not compile.
// One type's name must change.
type Constancy = {
  foo: boolean
}
type Constancy = {
  bar: boolean
}`,
			TestScript: `
// This koan is about the fact that type aliases cannot be redeclared.
// No runtime test needed.
`,
			TypeAssertions: `
// Only one type Constancy should exist
type _Check1 = Assert<IsType<Constancy, { foo: boolean } | { bar: boolean }>>;
`},
		{
			ID:          "interfaces",
			title:       "Interfaces",
			Label:       "",
			description: "An interface is very similar to a type",
			StarterCode: `??? MyInterface {
  foo: string;
  bar: number;
}`,
			TestScript: `
const obj: MyInterface = { foo: "hello", bar: 123 };
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
			StarterCode: `interface Person {
  name: string;
}
interface Monk ??? Person {
  isMeditating: boolean
}`,
			TestScript: `
const m: Monk = { name: "Huineng", isMeditating: true };
if (m.name !== "Huineng") throw new Error("Monk should have correct name");
if (!m.isMeditating) throw new Error("Monk should have isMeditating property true");
`,
			TypeAssertions: `
// Monk should extend Person and add isMeditating: boolean
type _Check = Assert<IsType<Monk, Person & { isMeditating: boolean }>>;
`},
		{
			ID:          "interfaces-redefining",
			title:       "Interfaces: Redefining",
			Label:       "",
			description: "An interface can be redefined freely",
			StarterCode: `interface MyInterface {
  foo: string;
}
??? MyInterface {
  bar: number;
}`,
			TestScript: `
const obj: MyInterface = { foo: "hi", bar: 5 };
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
			StarterCode: `function foo(myTuple: [string, ???]): true as const {
  return typeof myTuple[0] === "string"
  && typeof myTuple[1] === "number";
}`,
			TestScript: `
if (!foo(["a", 1])) throw new Error('foo(["a", 1]) should return true');
if (foo(["a", "b"] as any)) throw new Error('foo(["a", "b"]) should return false');
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
			StarterCode: `function foo(myTuple: ??? [string, number]): void {
  console.log(myTuple[0] + " will always be a string")
}`,
			TestScript: `
// Should accept a readonly tuple
const tuple: Readonly<[string, number]> = ["foo", 42] as const;
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
			description: "There exists a special type for functions that return promises",
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
			StarterCode: `type ManyThings = "one" | "another" | ???
let thing = "a secret third thing"`,
			TestScript: `
let t: ManyThings = "one";
t = "another";
t = "a secret third thing";
if (t !== "a secret third thing") throw new Error('t should be "a secret third thing"');
`,
			TypeAssertions: `
// ManyThings should include the string "a secret third thing"
type _Check = Assert<IsType<"a secret third thing", ManyThings>>;
`},

		{
			ID:          "literal-types-unions-of-numbers",
			title:       "Literal Types: Unions Of numbers",
			Label:       "",
			description: "A type can be a union of numbers",
			StarterCode: `type ManyNumbers = 1 | 2 | ???
let myNumber = 100`,
			TestScript: `
let n: ManyNumbers = 1;
n = 2;
n = 100;
if (n !== 100) throw new Error("n should be 100");
`,
			TypeAssertions: `
// ManyNumbers should include the number 100
type _Check = Assert<IsType<100, ManyNumbers>>;
`,
		},
		{
			ID:          "literal-types-as-literal",
			title:       "Literal Types: as Literal",
			Label:       "",
			description: "Literal types may require assertion",
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
// Colors.Blue should have the value 2
type _Check = Assert<IsType<typeof Colors.Blue, number>>;
`},
		{
			ID:          "enums-string",
			title:       "Enums: string",
			Label:       "",
			description: "Enums can have string values",
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
// Colors.Blue should have the value "BLUE"
type _Check = Assert<IsType<typeof Colors.Blue, string>>;
`},
		{
			ID:          "type-guards-typeof",
			title:       "Type Guards: typeof",
			Label:       "",
			description: "`typeof` can be used in expressions or in types",
			StarterCode: `let foo = "foo";
let bar: ??? foo;
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
			StarterCode: `type Person = {
  name: string
}
type Object = {
  foo: string
}
function typeDecider(thing: Person | Object): string {
  return "name" ??? thing ? "Person" : "Object";
}`,
			TestScript: `
if (typeDecider({ name: "Linji" }) !== "Person") throw new Error('typeDecider({ name: "Linji" }) should return "Person"');
if (typeDecider({ foo: "bar" }) !== "Object") throw new Error('typeDecider({ foo: "bar" }) should return "Object"');
`,
			TypeAssertions: `
// typeDecider should return "Person" or "Object" as a string
type _Check = Assert<IsType<ReturnType<typeof typeDecider>, string>>;
`},
		{
			ID:          "narrowing-instanceof",
			title:       "Narrowing: instanceof",
			Label:       "",
			description: "`instanceof` can be used to narrow types",
			StarterCode: `type Person = {
  name: string
}
type Object = {
  foo: string
}
function typeDecider(thing: Person | Object): string {
  return "name" ??? Person ? "Person" : "Object";
}`,
			TestScript: `
// instanceof can't be used directly with plain object types, so let's provide classes:
class PersonClass { constructor(public name: string) {} }
class ObjectClass { constructor(public foo: string) {} }

function classTypeDecider(thing: PersonClass | ObjectClass): string {
  return thing instanceof PersonClass ? "Person" : "Object";
}

if (classTypeDecider(new PersonClass("Linji")) !== "Person") throw new Error('classTypeDecider(new PersonClass("Linji")) should return "Person"');
if (classTypeDecider(new ObjectClass("bar")) !== "Object") throw new Error('classTypeDecider(new ObjectClass("bar")) should return "Object"');
`,
			TypeAssertions: `
// classTypeDecider should return "Person" or "Object" as a string
type _Check = Assert<IsType<ReturnType<typeof classTypeDecider>, string>>;
`},
		{
			ID:          "type-predicates-is",
			title:       "Type Predicates: is",
			Label:       "",
			description: "A type predicate will tell tthe compiler about the type of a variable",
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
			description: `Index signatures let you type objects with unknown keys, but known value types.`,
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
			description: "The unknown type is a safer alternative to any.",
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
type _AssertName = Assert<IsType<typeof user, HasName>>;
type _AssertAge = Assert<IsType<typeof user, HasAge>>;
`,
		},
		{
			ID:          "generics-type-alias",
			title:       "Generics: Type Alias",
			Label:       "",
			description: `You can create reusable types with generics.`,
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
			ID:          "keyof",
			title:       "The keyof Keyword",
			Label:       "",
			description: `keyof returns a union of `,
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
			StarterCode: `type MaybeUser = {
    id?: number;
    username?: string;
    email?: string;
}

type RequiredUser = {
    ??? [K in keyof MaybeUser]: MaybeUser[K]
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
// Should not be readonly; use -readonly tto remove the readonly property
type _Check = Assert<IsNotType<WritableUser["id"], Readonly<number>>>;
`,
		},
		{
			ID:          "utility-types-partial",
			title:       "Utility Types: `Partial`",
			Label:       "",
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
			StarterCode: `type User = {
    id?: number;
    username?: string;
    email?: string;
}

type FullUser = ???

const u: FullUser = {
    id: 100,
    username: "Yun-men",
    email: "yun-men@suemru.com"
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
			title:       "Utility Types: Omit`",
			Label:       "",
			description: "Sometimes you must omit, to create something new",
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
if (user.email !== "Dongshan@example.com") throw new Error("email should be 'Dongshan@shouchu.com'");
`,
			TypeAssertions: `
// Should not allow password
type _Check = Assert<IsNotType<keyof PublicUser, "password">>;
type _Check2 = Assert<IsType<keyof PublicUser, "id" | "username" | "email">>;
`,
		},
	}

	return exercises
}
