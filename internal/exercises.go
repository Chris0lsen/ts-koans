package internal

type Exercise struct {
	ID          string
	title       string
	description string
	StarterCode string
	TestScript  string
	Label       string
	FunctionName string
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
        ID: "add",
        title: "Add Two Numbers",
        description: "Write a function that adds two numbers. Both arguments and the return value must be numbers.",
        StarterCode: `
export function add(a, b) {
  // TODO: Add two numbers
}
`,
        TestScript: `
if (typeof sandbox.exports.add !== "function") throw new Error("Not a function");
if (sandbox.exports.add(2, 3) !== 5) throw new Error("add(2, 3) !== 5");
console.log("✅ All tests passed!");
`,
        Label: "",
        FunctionName: "add",
        TypeAssertions: `
// Type check errors below mean your function signature or return type is not correct!
// Fix your function until no errors appear here.

type Assert<T extends true> = T;

// Assert first parameter is a number
type AssertAddParam0IsNumber = Assert<Parameters<typeof add>[0] extends number ? true : false>;

// Assert second parameter is a number
type AssertAddParam1IsNumber = Assert<Parameters<typeof add>[1] extends number ? true : false>;

// Assert return type is number
type AssertAddReturnType = Assert<ReturnType<typeof add> extends number ? true : false>;
`,
    },
    {
        ID: "make-user",
        title: "Make User Object",
        description: "Create a function that returns a user object with properties 'name' (string) and 'age' (number).",
        StarterCode: `
export function makeUser(name, age) {
  // TODO: Return an object with 'name' and 'age'
}
`,
        TestScript: `
const u = sandbox.exports.makeUser("Alex", 42);
if (!u || u.name !== "Alex" || u.age !== 42) throw new Error("Incorrect user object");
console.log("✅ All tests passed!");
`,
        Label: "",
        FunctionName: "makeUser",
        TypeAssertions: `
// Type check errors below mean your function signature or return type is not correct!
// Fix your function until no errors appear here.

type Assert<T extends true> = T;
type MustBeUser = { name: string; age: number; };

// Assert first parameter is string
type AssertMakeUserParam0IsString = Assert<Parameters<typeof makeUser>[0] extends string ? true : false>;

// Assert second parameter is number
type AssertMakeUserParam1IsNumber = Assert<Parameters<typeof makeUser>[1] extends number ? true : false>;

// Assert return type matches { name: string; age: number }
type AssertMakeUserReturnType = Assert<ReturnType<typeof makeUser> extends MustBeUser ? true : false>;
`,
    },
    {
        ID: "optional-nickname",
        title: "Optional Properties",
        description: "Add an optional property 'nickname' (string) to the User interface and implement getNickname.",
        StarterCode: `
interface User {
  name: string;
  // TODO: Add an optional property 'nickname'
}

export function getNickname(user: User): string | undefined {
  // TODO: Implement me!
}
`,
        TestScript: `
if (sandbox.exports.getNickname({ name: "A" }) !== undefined) throw new Error("Should return undefined for missing nickname");
if (sandbox.exports.getNickname({ name: "A", nickname: "B" }) !== "B") throw new Error("Should return nickname");
console.log("✅ All tests passed!");
`,
        Label: "",
        FunctionName: "getNickname",
        TypeAssertions: `
// Type check errors below mean your function signature or return type is not correct!
// Fix your function until no errors appear here.

type Assert<T extends true> = T;
type HasOptional<T, K extends keyof T> =
  {} extends Pick<T, K> ? true : false;
type UserParam = Parameters<typeof getNickname>[0];

// Assert that 'nickname' is optional on User
type AssertNicknameIsOptional = Assert<HasOptional<UserParam, "nickname">>;

// Assert return type
type AssertGetNicknameReturnType = Assert<ReturnType<typeof getNickname> extends string | undefined ? true : false>;
`,
    },
    {
        ID: "union-type",
        title: "Union Types",
        description: "Write a function that accepts a string or number and returns its string representation.",
        StarterCode: `
export function toStringValue(x) {
  // TODO: Implement me!
}
`,
        TestScript: `
if (sandbox.exports.toStringValue(42) !== "42") throw new Error("Should stringify numbers");
if (sandbox.exports.toStringValue("hi") !== "hi") throw new Error("Should return strings");
console.log("✅ All tests passed!");
`,
        Label: "",
        FunctionName: "toStringValue",
        TypeAssertions: `
// Type check errors below mean your function signature or return type is not correct!
// Fix your function until no errors appear here.

type Assert<T extends true> = T;

// Assert parameter is string or number
type AssertToStringValueParam = Assert<Parameters<typeof toStringValue>[0] extends string | number ? true : false>;

// Assert return type is string
type AssertToStringValueReturnType = Assert<ReturnType<typeof toStringValue> extends string ? true : false>;
`,
    },
    {
        ID: "keyof-type",
        title: "keyof Types",
        description: "Write a function that takes an object and a key, and returns the value at that key. Use keyof for type safety.",
        StarterCode: `
export function getValue(obj, key) {
  // TODO: Implement me!
}
`,
        TestScript: `
const obj = { a: 1, b: 2 };
if (sandbox.exports.getValue(obj, "a") !== 1) throw new Error("getValue did not work for 'a'");
if (sandbox.exports.getValue(obj, "b") !== 2) throw new Error("getValue did not work for 'b'");
console.log("✅ All tests passed!");
`,
        Label: "",
        FunctionName: "getValue",
        TypeAssertions: `
// Type check errors below mean your function signature or return type is not correct!
// Fix your function until no errors appear here.

type Assert<T extends true> = T;
type O = { a: number; b: number; };

// Assert first parameter is O
type AssertGetValueParam0 = Assert<Parameters<typeof getValue>[0] extends O ? true : false>;

// Assert second parameter is keyof O
type AssertGetValueParam1 = Assert<Parameters<typeof getValue>[1] extends keyof O ? true : false>;

// Assert return type is number
type AssertGetValueReturnType = Assert<ReturnType<typeof getValue> extends number ? true : false>;
`,
    },
    {
        ID: "readonly-array",
        title: "Readonly Arrays",
        description: "Write a function that accepts a readonly array of numbers and returns their sum.",
        StarterCode: `
export function sumReadonly(nums) {
  // TODO: Implement me!
}
`,
        TestScript: `
if (sandbox.exports.sumReadonly([1,2,3]) !== 6) throw new Error("Sum incorrect");
console.log("✅ All tests passed!");
`,
        Label: "",
        FunctionName: "sumReadonly",
        TypeAssertions: `
// Type check errors below mean your function signature or return type is not correct!
// Fix your function until no errors appear here.

type Assert<T extends true> = T;

// Assert parameter is readonly number[]
type AssertSumReadonlyParam = Assert<Parameters<typeof sumReadonly>[0] extends readonly number[] ? true : false>;

// Assert return type is number
type AssertSumReadonlyReturnType = Assert<ReturnType<typeof sumReadonly> extends number ? true : false>;
`,
    },
    {
        ID: "typeof-inference",
        title: "typeof Inference",
        description: "Use 'typeof' to type a new variable after a constant has been declared.",
        StarterCode: `
const ANSWER = 42;

// TODO: Type x using 'typeof'
const x = ANSWER;
export { x }
`,
        TestScript: `
if (sandbox.exports.x !== 42) throw new Error("x should be 42");
console.log("✅ All tests passed!");
`,
        Label: "",
        FunctionName: "x",
        TypeAssertions: `
// Type check errors below mean your function signature or return type is not correct!
// Fix your function until no errors appear here.

type Assert<T extends true> = T;

// Assert x has the same type as ANSWER
declare const ANSWER: 42;
declare const x: typeof ANSWER;
type AssertTypeofX = Assert<typeof x extends typeof ANSWER ? true : false>;
`,
    },
}


	return exercises
}
