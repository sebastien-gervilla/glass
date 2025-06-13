# Glass

Glass is an interpreted language created in Golang.
It's an attempt at learning more about languages in general, and therefore how to create them.

NOTE: The language is still in progress, so there surely are missing features and bugs :)

## Testing

First Glass needs to be built :

`go build ./cli/main.go`

Then you can create a glass file in the project's root :

`./glass/main.glass`

As an example, paste :

```
let MAXIMUM_AGE = 100;

let isAgeValid = fn(age) {
    if (age < 0) {
        return false;
    };

    return age < MAXIMUM_AGE;
}

let age = 18;
if (isAgeValid(age)) {
    print(age, " is a valid age");
} else {
    print(age, " is not a valid age");
};
```

Finally, it can be executed with the command :

`./main.exe run ./glass/main.glass`

## Features

It mostly support basic features such as :

### Variable definitions

`let x = 5;`

### Basic operations

`let x = (5 + 6) * 4;`

### Conditions

```
if (x > 30) {
    print("x greater than to 30")
} else {
    print("x inferior or equal to 30")
};
```

### Functions

Functions are declared as variables.

```
let add = fn(a, b) {
    return a + b;
}
```

### Builtins

You can log into the console by using `print` :

`print("x is equal to ", x);`