# Test Driven Design Steps

These are the steps from the books:
* The Deeper Love Of Go
* The Power Of Go: Tools
* The Power Of Go: Tests

They are distilled for clarity.

## Decide Scope

* Start by trying to describe the behaviour that we want in words.
* Refine that description into the form of a test.

## Creating Tests

* Create a `APP_test.go` file.
* Start with the function you are testing. It doesn't have to exist right now
* Work backwards to create the necessary settings for the test to run.
* Work forwards to check the results against expectations.
* Once the test is written, run it so you see it fails.
* Now write a null implementation of the function you are testing. Re-running the test should fail again.
* Write the real implementation of the function you are testing. Run the test and fix the function code until it succeeds.

## Test Files

* Name the files after the test case they represent.
  * e.g. if we are counting lines, named the file `3_lines.txt` and place 3 lines of text into it.

## Red, Green, Refactor

Write your test as though there was a "magic" package that does whatever you want.

* It takes the inputs you have,
* It returns what you want,
* All paperwork is dont inside the "magic" package.
* Implement the "magic" package.
  * Tests fail first (red).
  * Tests pass next (green).
  * Refactor the "magic" package code so it is easier to read and maintain. The tests should still pass.

## ACE (Action, Condition, Expectation)

Test names should be ACE.

* The **Action** performed.
* The **Conditions** being tested.
* The **Expected** results.

## Testing Questions

* What are we really testing here?
* What intent is this test communicating?

## Summary Of Creating, Inspecting, & Testing Errors

1. Always check errors in tests.
2. We usually don't care what type an error is, as long as they aren't nil.
3. When it does matter what type an error is, use `errors.Is` to compare them with a sentinel error.
4. When combining a sentinel error with dynamic information, use `fmt.Errorf` with the `%w` verb to create a wrapped error.

## Generating Useful Test Inputs

* Ask the question: what could go wrong here?
* Choose inputs from each equivalence class. e.g. positive numbers, negative numbers, integers, floats, 1, 0, -1.
* `uint` can be used to force positive numbers, but it is better to test inputs being >= 1 instead of relying on this.
* For slices, pass in a nil slice or empty slice.
* For pointers, pass in nil.

## Table Tests

* A table test performs the same check on each set of test cases, which are typically stored inside a struct, but a map can be used with t.Run subtests to do this elegantly with the subtest name stored as the key.
* Have seperate tests for valid input and invalid input.

```go
testCases := map[string]struct {
input []string
want  string
}{
  // The map's key is often the name of the test used by the subtest function.
  "no items": {
    input: []string{},
    want:  "",
  },
}
for name, tc := range testCases {
  t.Run(name, func(t *testing.T) {
    got := game.ListItems(tc.input)
    if tc.want != got {
      t.Error(cmp.Diff(tc.want, got))
    }
  })
}
```

## Test Data

* Use fake data when applicable, so readers of the test case know that data is irrelevant to the test.

```go
got, err := CallAPI(request, "fake API token")
```

* Rather than using global variables in the test package, create a function that returns the data. This protects against a test modifying test data and producing weird results on other tests. Use this for input that contains maps and slices, including structs whose fields are maps or slices.

```go
func makeAgeData() map[string]int {
  return map[string]int{
    "alice": 18,
    "bob": 99,
  }
}
```

* If test data gets too big, store it in a file inside of [test/data](https://github.com/golang-standards/project-layout#test), but files are slow so use alternatives where possible.
* Try using io.Reader/Writer instead of files, which can be handled easily by a bytes.Buffer. eg strings.NewReader("some string") for io.Reader and io.Discard or bytes.Buffer for io.Writer.
* Use fs.FS for trees of files on disk. See chapter 5 of The Power Of Go: Tools.
* Use t.TempDir for output with cleanup, combine with t.Name to get the name of the test. `t.TempDir() + "/" + t.Name() + ".txt")
* 
