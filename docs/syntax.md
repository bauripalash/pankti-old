## [DRAFT] Pankti Programming Language : Syntax 

**Status**
	This document specifies the syntax of the pankti programming 
	language and requests discussion and suggestion for improving the 
	syntax and usability of the aforementioned programming language.

**Copyright Notice**
	Copyright (C) Palash Bauri. All Rights Reserved

**Abstract**
	Pankti Programming language is a specially designed programming 
	language to be used with the Bengali language. It uses hand written
	lexer and Pratt parser. This document will specify and act as an 
	reference point for all the implementation of Pankti Programming 
	Language.

**Table of Contents**
* Introduction 
* Assignments 

### Introduction
Pankti programming language is an Bengali programming language is 
designed to easy to write as well as easy to read. As the language 
targets technical and non-technical population, it must be simple and 
all the required complexities of programming and computer science must 
be presented in easy to understand way.


### Definition
The name of the programming language is **Pankti Programming Language**
. In short it can be mentioned as **pankti**. The default 
implementation of **pankti** is written in Go and should be called as
**pank**; The first prototype of pankti would be mentioned as **P0**
here and the current new syntax version of pankti would be mentioned as 
**P1**


### Basic Syntax

At the first prototype of pankti had many different syntax than current
iteration of pankti. 

**P0** Used to have this type of `IF..ELSE` syntax
```go
jodi (true) tahole{
    return true
}
```

```go
jodi (true) tahole{
    return true
}else{
    return false
}
```

But now **P1** has this type of `IF..ELSE` syntax
```go
jodi (true) tahole 
    return true
nahole sesh
```

```go
jodi (true) tahole 
    return true
nahole 
    return false
sesh
```
