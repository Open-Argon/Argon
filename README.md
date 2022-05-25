# CURRENTLY IN DEVELOPMENT!

<div align="center">
<p>
    <img width="80" src="https://raw.githubusercontent.com/Ugric/Argon/main/logo.png">
</p>
<h1>The Argon Programming Language</h1>
</div>

## Key Features of Argon

- easy to learn for new developers
- automatic type changing (e.g. INPUT: `'age: ' + 34` OUTPUT: `'age: 34'`)
- boolean renamed from true and false to yes and no
- the syntax is very similar to python and javascript
- no floats and ints, only numbers
- **New Feature!** easier operators, e.g. `5^2` can now be `5 to the power of 2`

## Examples

Hello world:

```javascript
log("hello, world!")
```

Inputs:

```javascript
input("name: ")
```

Variables:

```javascript
const myConstant = "this can't change!"
var myVariable = "this can change!"
myOtherVariable = "this can also change!"
```

Addition:

```javascript
const word1 = "hello"
const word2 = "world"
log(word1 + " " + word2)
```

If Statements:

```javascript
const password = 'password123'
const guess = input('password: ')
if (password == guess) [
   log('correct')
] else [
   log('incorrect')
]
```

While Loops:

```javascript
var i = 0
while (i < 100) [
   log(i)
   i++
]
```

Logic (Boolean):

```javascript
log("true is now", yes)
log("false is now", no)
log("none is now", unknown)

while (yes) [
   log("Argon is EPIC!")
]
```

Items (Lists or Arrays):

```javascript
[1, 2, 3]
const names = ["john", "joe", "bob"]

log(names[2])
```

Books (Dictionaries or Objects):

```javascript
const user = {'name': "bob", 'age': 34, 'hobbies': ["Programming in Argon!", "Playing video games!"]}
log(user["name"])
```

Subs (Functions):

```javascript
sub hello(username) [
  return 'hello, ' + username + ', how are you?'
]
log(hello('bob'))
```

Easy to Learn Operators:
```javascript
log('addtion:')
log(1 + 2)
log(1 add 2)
log(1 plus 2)
log()
log('subtraction:')
log(1 - 2)
log(1 subtract 2)
log(1 minus 2)
log()
log('multiplication:')
log(1 * 2)
log(1 multiplied by 2)
log(1 times 2)
log(1 x 2)
log()
log('division:')
log(1 / 2)
log(1 divided by 2)
log(1 over 2)
log()
log('modulo:')
log(1 % 2)
log(1 modulo 2)
log(1 mod 2)
log()
log('power:')
log(1 ^ 2)
log(1 to the power of 2)
log(1 ** 2)
log()
log('equality:')
log(1 == 2)
log(1 is 2)
log()
log('inequality:')
log(1 != 2)
log(1 is not 2)
log()
log('greater than:')
log(1 > 2)
log(1 is greater than 2)
log()
log('less than:')
log(1 < 2)
log(1 is less than 2)
log()
log('greater than or equal to:')
log(1 >= 2)
log(1 is greater than or equal to 2)
log()
log('less than or equal to:')
log(1 <= 2)
log(1 is less than or equal to 2)
log()
log('and:')
log(1 && 2)
log(1 and 2)
log()
log('or:')
log(1 || 2)
log(1 or 2)
```

## How do I run Argon?

Setup:

- download the Argon executable for [Windows](https://github.com/Ugric/Argon/raw/main/dist/Windows/argon.exe), [MacOS](https://github.com/Ugric/Argon/raw/main/dist/macOS/argon) or [Linux](https://github.com/Ugric/Argon/raw/main/dist/Linux/argon)

Run a file:

- make a file with the extention as `.ar` and put your Argon code inside that file
- call the Argon executable with a parameter as the path to your `.ar` file, (e.g. `$ argon example`)

Run the shell:

- run the Argon executable without any parameters, (e.g. `$ argon`)
