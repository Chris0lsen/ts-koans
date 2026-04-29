// -- Type Assertion Utilities --

// Produces a type error if 'T' is not 'true'
type Assert<T extends true> = T;

// Checks if two types are the same (structurally)
// https://stackoverflow.com/questions/53807517/how-to-test-if-two-types-are-exactly-the-same
type IsType<A, B> = (<T>() => T extends A ? 1 : 2) extends
                    (<T>() => T extends B ? 1 : 2) ? true : false;

// Checks if type A is not assignable to B
type IsNotType<A, B> = IsType<A, B> extends true ? false : true;

type IsNotReadonly<T, K extends keyof T> =
  IsType<Pick<T, K>, { -readonly [P in K]: T[P] }>;


// Checks if type A is assignable to B
type IsAssignable<A, B> = A extends B ? true : false;
