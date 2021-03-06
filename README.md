# zen-lang

Zen is a proof of concept language for compile time preconditions, postconditions, and invariants

Consider the example below

```golang
func Min[a: i32, b: i32] => [result: i32]
    where result <= a && result <= b
{
    if (a < b) {
        return a
    } else {
        return b
    }
}
```

The compiler is aware of the constraint `result <= a && result <= b` so if a bug were introduced by changing `a < b` to `a > b` the compiler would output

```
Could not verify post conditions
With precondition at
../../test/Min.zen: (2, 26)
    where result <= a && result <= b
                         ^
../../test/Min.zen: (5, 9)
        return a
        ^
```

Invariants can be applied to data structures too

```golang
type Range [
    Min: i32,
    Max: i32,
] where Min <= Max
```

With this, the compiler can catch bugs around malformed data structures instead of relying on runtime assertions. So this would be a compile time error

```golang
var range: Range = [2, 1]
```

Creating a well formed range from any two integers could be implemented with

```golang
func MakeRange[a: i32, b: i32] => [result: Range] {
    if (a < b) { // the compiler takes note of this comparison here
        return [a, b] // to make sure a < b here
    } else {
        return [b, a] // and a >= b here
    }
}
```

Using compile time constraints the complier can prevent errors such as out of bounds exceptions and divide by zero errors with the goal of eliminating all runtimes errors normally caught by assertions and making them compile time errors

For more details on how the algorithm works see [Bounds Checking](./doc/boundschecking.md)