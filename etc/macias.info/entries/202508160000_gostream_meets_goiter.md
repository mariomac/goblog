# Go streams meet standard Go iterators!

In a previous blog post, titled
[Performance comparison of Go functional stream libraries](/entry/202212020000_go_streams.md),
we compared the performance and features of several Go functional stream libraries.
Since then, the Go standard library has added the [iter package](https://pkg.go.dev/iter),
which defines a foundation for iterators, as well as many other packages allowing
to manipulate [slices](https://pkg.go.dev/slices) and [maps](https://pkg.go.dev/maps)
in a functional style.

After the new additions to the standard library, most programmers might feel
discouraged from using a third-party stream library if the standard library packages
already provide a concise way of manipulating slices and maps.

In an effort to keep the [gostream library](https://github.com/mariomac/gostream)
up to date with the latest Go idioms,
the version 0.10.1 of [gostream](https://github.com/mariomac/gostream) introduced
new functions and methods to let it work together with the [iter package](https://pkg.go.dev/iter),
and the rest of the packages using it.

## The core of the `iter` package: `iter.Seq`

The Go iter package defines the `iter.Seq` generic type as:

```go
type Seq[V any] func(yield func(V) bool)
```

Providing an iterator to your custom data structures means providing
a function that will invoke the `yield` function for each element of that
structure. The following `BackwardsCounter`
type shows an example about how to provide an iterator that would iterate
backwards all the numbers between a given value and 1:

```go
type BackwardsCounter int

func (b BackwardsCounter) Iterator() iter.Seq[int] {
    return func(yield func(int) bool) {
        for count := int(b) ; count > 0 ; count-- {
            if !yield(count) {
                return
            }
        }
    }
}
```

Then your `BackwardsCounter.Iterator()` method can be used in multiple standard
functions using `iter.Seq`. Even in a `for ... range` clause:

```go
func main() {
    countDown := BackwardsCounter(5)
    for n := range countDown.Iterator() {
        fmt.Println(n)
    }
}
```

For map-like data structures formed by sequences of key-values, the generic
`iter.Seq2[K, V]` behaves similarly but its `yield` function argument now
takes two values instead of one: the Key and Value of each mapped entry.

```go
func AnimalSounds() iter.Seq2[string, string] {
    return func(yield func(string, string) bool) {
        yield("Dog", "wow wow!")
        yield("Cat", "meeeow")
        yield("Frog", "ribbit!")
    }
}

func main() {
    for k, v := range AnimalSounds() {
        fmt.Println(k, "says", v)
    }
    fmt.Println("All together:", maps.Collect(AnimalSounds()))
}
```

Output:

```text
Dog says wow wow!
Cat says meeeow
Frog says ribbit!
All together: map[Cat:meeeow Dog:wow wow! Frog:ribbit!]
```

## Why keep using `gostream`?

Despite the powerful additions to the Go standard library, `gostream` still
provides unique advantages including lazy evaluation and support for
infinite streams, along with additional convenience methods not yet covered
by the standard library (for example, map-reduce operations).

However, do not use `gostream` just because it's elegant and cool 😎, if you
don't feel it might provide any benefit that justifies adding a new dependency
to your code.

## Connecting `gostream` library with the `iter` package

### Instantiation of a `stream.Stream` from an `iter.Seq`

The `stream.OfSeq` function allows building a `stream.Stream[T]` from a given
`iter.Seq[T]`:

```go
nums := stream.OfSeq(func(yield func(string) bool) {
    yield("one")
    yield("two")
    yield("three")
}).Map(strings.ToUpper).ToSlice()

fmt.Println(nums)
```

Output:

```text
[ONE TWO THREE]
```

The `stream.OfSeq2` function will take an `iter.Seq2[K, V]` as input and will
create a `stream.Stream[item.Pair[K, V]]`:

```go
stream.OfSeq2(func(yield func(int, string) bool) {
    yield(1, "one")
    yield(2, "two")
    yield(3, "three")
}).ForEach(func(i item.Pair[int, string]) {
    fmt.Println(i.Val, "is", i.Key)
})
```

Output:

```text
one is 1
two is 2
three is 3
```

### How to iterate a `stream.Stream[T]`

Before integrating [gostream](https://github.com/mariomac/gostream) with the
standard `iter` package, the only way to iterate a `Stream` was to use the
`.ForEach` method (in addition to other collection methods that internally
iterate the streams, for example `.ToSlice`, `.Count`, `.AnyMatch` and so on).

The latest version of `gostream` provides the following helper functions for iterating
a stream.

#### Iterating sequential streams

You can iterate a sequential stream inside a `for ... range` loop in two ways:
* The `.Seq()` method returns an `iter.Seq` that iterates over each element
  of the `Stream`, in order.
* The `.Iter()` method returns an `iter.Seq2` iterator where the key type is the index
  of the element within the sequence, and the value type is the element itself.

For example, given a `oneToFive` stream containing the values `1, 2, 3, 4, 5`:

```go
oneToFive := stream.Iterate(1, item.Increment[int]).
	Limit(5)

for val := range oneToFive.Seq() {
    fmt.Print(val, " ")
}
fmt.Println()

for index, val := range oneToFive.Iter() {
    fmt.Printf("oneToFive[%v] = %v\n", index, val)
}
```

Output:

```text
1 2 3 4 5 
oneToFive[0] = 1
oneToFive[1] = 2
oneToFive[2] = 3
oneToFive[3] = 4
oneToFive[4] = 5
```

#### Iterating key-value streams

If your stream is: `stream.Stream[item.Pair[K, V]]`, you can pass it to the
`stream.Seq2` function so you can iterate it in a `for ... range` loop via an
iterator where the key element is the `item.Pair` key type `K` and the value
element is the value type `V`:

```go
petSounds := stream.OfMap(map[string]string{
    "dog": "woof",
    "cat": "meow",
    "bird": "chirp",
    "fish": "blub",
})

for pet, sound := range stream.Seq2(petSounds) {
    fmt.Println( pet, "says", sound)
}
```

Output:

```text
dog says woof
cat says meow
bird says chirp
fish says blub
```

## Connecting `gostream` with the rest of the Go standard library

The [slices](https://pkg.go.dev/slices) and [maps](https://pkg.go.dev/maps)
standard Go packages provide functions either accepting or returning both
`iter.Seq` and `iter.Seq2`.

You can use the aforementioned `gostream` methods
for instantiating streams from standard iterators (`stream.OfSeq`,
`stream.OfSeq2`) to connect the outputs of the Go standard libraries as inputs
to `gostream`.

Inversely, you can use the `gostream` methods and functions for creating 
standard iterators (`Iter()`, `Seq()`, `Seq2(...)`) to connect the output
of `gostream` to the Go standard library.

For example:

```go
animalSounds := map[string]string{"dog": "woof", "cat": "meow", "cow": "moo"}

// returns an iter.Seq[string]
animals := maps.Keys(animalSounds)

// 1. transform the iterator to a stream
// 2. apply the strings.ToUpper function to each element
// 3. convert the stream back to an iter.Seq
upper := stream.OfSeq(animals).
    Map(strings.ToUpper).
    Seq()

// keep using the data flow with the standard library
fmt.Println(slices.Sorted(upper))
```

Output:

```text
[CAT COW DOG]
```

Check the documentation of the [slices](https://pkg.go.dev/slices) and
[maps](https://pkg.go.dev/maps) functions for more examples of standard functions
accepting and returning iterators.



