
### var declaration

var name ***type*** = ***expresion***
either the type or the expression can be ommited, but not both

``` go
var i, j, k int 
var l, m n = 1, "2", true
var f, err = os.Open("/tmp/")
```

### short variable declaration

***name*** := ***expression***

``` go
i := 0
j := "3"
integers := []int {1,2,3}
```

### := and =

:= is a declaration, whereas = is a assignment
``` go
i, j := 1, 3 // declaration
i, j = j, i // swap/assignment
```

### pointer and address

***&x***: the address of x
****int***: pointer to int

``` go
x := 1
p = &x // assign address of variable x to p, p is of type *int and point to x
*p = 2 // equivalent to x  = 2
```

### new()

``` go
p = new(int) // create a unnamed variable of type int, initialize it to 0 and return its address, which is a value of type *int
*p = 2
```

a variable created with new is no different from an ordinary local variable whose address is taken, except that there's no need to declare a dummy name

below have identical behaviors

``` go
func newInt() *int {
    return new(int)
}

```

``` go
func newInt() *int {
    var dummy int
    return &dummy
}

```

### type Declarations

it defines a new named type that has the same underlying tye as an existing type

> type name underlying-type

``` go
type Celsius float64
type Age  int

```

### variable shadow

``` go
func f(){}

func main(){
    f := "shadow" //local var f shadows the package-level var func f
}
```

### Strings

``` go
s := "hello world"
len(s)
s[0]
s[5]
s[:5] // hello
s[7:] // world
s[:] // hello, world
s[0] = '3' // error, string is immutable
```

### String to int and vise versa

``` go
s := strconv.Itoa(123) // integer to ascii
i := strconv.Atoi("123)
```

### const

``` go
const PI = 3.14
const (
    e = 2.7
    pi = 3.14
)
```

### constant generator: iota

iota default to 0

``` go
type Weekday int

const (
    Sunday Weekday iota
    Monday
    Tuesday
    Wednesday
    Thuresay
    Friday
    Saturday
)

type Flags uint
const (
    FlagUp Flags =  1 << iota
    FlagBroadcast
    FlagLookback
    FlagMulticast
)
```

### array

fixed-length

```
var a [3]int
var b [3]int = [3]int{1,4}
c :=[...]int{1,2,3} // determined by the number of initializers

for i, v :=range a {
    fmt.Print("%d %d \n", i, v)
}

```
### Slice
``` go
s := []int{1,2,3,4,5} // differs from array, where the [] is empty
s[0] = 0 // in-place update
len(s) == 0 // is empty

var ret []int
for _, r := range [4]int{6,7,8,9} {
    ret = append(ret, r)
}

ret = append(ret, 1,2,3,4)
ret = append(ret, ret)

```

### Map

``` go
ages := make(map[string]int) // mapping from strings to ints

ages := map[string]int {
    "alice": 12,
    "bob": 34,
}
ages["bob"]= 13
delete(ages,"bob")

from name, age: range ages {
    fmt.Print("%s %d \n", name, age)
}

age, ok := ages["bob"]

if ok {/*bob is in the map*/}

```

### Structs

``` go

type Employee struct {
    Name string
    Age int
    Address string
}

var employee Employee
// or
employee := Employee{
    Name: "bob",
    Age: 12
}

employee.Name = "bob"
employee.Age = 23


func Scale (e Employee) {
    e.Name = ""
    e.Age = 0
}
// larger struct types are usually passed to or returned from functions using a pointer
func Scale (e *Employee){
    // ....
}

//obtain address
p := &Employee{
    Name: "bob",
    Age: 12
}
//equals to

p := new(Employee)
*p =  Employee{
    Name: "bob",
    Age: 12
}


```

### functions

optional parameter list
optional result list

```
func name(parameter-list) (result-list){
    body
}
```

```go

func add (x int, y int) int {return x + y}
func add (x, y int) (z int) {z = x + y; return}
func Parse(r io.Reader) (*Node, error) {/*....*/}

func squares(x int) func() int { // function as type
    
    return func() int{
        return x * x
    }
}

func opt (f func( x int) int, y int) int {
    return f(y)
}

//variadic functions

func sum (values ...int) int {
    total := 0
    return for _, v : range values {
        total += v
    }

    return total
}

```
### defer and recover

``` go
func A(){

    /*
     .
     .
     .
    */    
   defer f.Close(); // execute until the function A has finished.
   return list;
}


// https://www.cnblogs.com/charlieroro/archive/2018/03/21/8617056.htm
func Parse(input string) (s *Syntax, err error) {
   
    defer func(){ // execute when parse finished, if there are panic during the parsing, 
        if p:= recover(); p!=nil {
            err = fmt.Error("internal error: %v", p)
        }
    }

    // parsing logic
}

```