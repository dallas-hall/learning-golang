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
