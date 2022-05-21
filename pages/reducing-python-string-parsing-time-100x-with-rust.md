---
title: Reducing python string parsing time 100x with rust
date: "2021-04-25T22:40:32.169Z"
---

*april 2021*


Recently my work focused on the development of IBM's LSF simulator. LSF is a workload scheduler for HPC and one of the
problems we were aiming to solve is how to simulate the running of large HPC clusters. Like many products, especially
those who's inception was 15+ years ago, a custom logging format is used.  These log files can have millions of
daily job entries and parsing these files was a requirement before simulations could start.


The current parser was pure
python3 and could take a significant amount of time when the number of entries was large.Using GDB it was possible to
determine that the bottleneck were python operations used in the actual parsing of strings, namely the large amount of
string copies as they are immutable in the language it was written in and thus concatenation and modification of strings
can require many allocations. Room did exists to use, "better" python primitives for the operations required, however for
the desired speed up I selected to create a python extension that did all parsing in rust/c. I created a proof of concept
in rust, however for better integration with current build tooling a c implementation was selected for production.

What was surprising was how easy creating a python extension in rust was and how easy it was to generate correct code vs
the final c implementation (the correctness should come at no surprise to anyone familiar with rust with it being slightly higher level as well as
the power of it's compiler)

Typically, as discussed above, a python extension is written in C, however through the [pyo3](https://github.com/PyO3/pyo3) we can write
our extensions in rust and work with python objects directly as py03 acts a DSL that compiles directly to cPython
(see [Ben Frederickson](https://www.benfrederickson.com/writing-python-extensions-in-rust-using-pyo3/) excellent blog).

### Why Rust ?

- No GC, to avoid the slowdown seen above I wanted to avoid any Garbage Collected Language. This excludes go, another
  favourite of mine as despite the reduced memory footprint of go escaping the GC would be critical to significant speed-ups.

- Heap and Stack control, having to essentially explicitly declare where memory is allocated provided more opportunity for
finding optimizations as well as forces more in depth thought on how we were storing and passing state as well.

### Why Rust over C Initially?

- Memory safety, I had increased confidence in my ability to deliver run-time error free code in rust and wanted to be
  freed from worrying about memory leaks.

- Higher level string abstractions, though c string libraries exist being able to use String or &str over char**  with no
  dependencies is a blessing (c dependencies outside of the norm where frowned upon). Correctness here is still valued
  over any speed up the use of raw pointers could provide.

- Development speed, delivering a correct, leak free program would take significantly less time in Rust vs C.

Simply rust provided a excellent balance between performance and safety as well as a excellent development experience.
### pyO3

Using python as a high level interface for low level execution is not a new concept, most of python's standard library
exists as a collection of c extension. As well as many popular libraries such Pytorch (which now uses rust), Tensorflow,
Numpy. PyO3 makes it dead simple to create python modules in pure rust. When using any extension for performance gains
the things in mind are how you are interacting with the GIL and how Python Objects are stored/returned. PyO3 abstracts
away a ton of the complexity for us as rust to python type conversion is automatic natively while retaining the ability
to explicitly work on python objects themselves. When working in Python extensions in c, yes ctypes can be returned easily,
but python objects themselves which includes the string class mst be explicitly created.

### A parsing experiment

For an example, we will attempt to create a python extension in rust to allow us to split strings.

The below code is a replacement for pythons str.split(), with a few caveats...

1. only supports splitting on a single char
2. only accepts US-ASCII characters, i.e 1 byte is 1 char allowing us to operate directly on the bytes

```
use pyo3::prelude::*;
use pyo3::wrap_pyfunction;

#[pyfunction]
pub fn split(s: &str , character: char) -> Vec<&str> {

    let mut current_start;
    let mut current_end = 0;
    let str_length = s.chars().count();
    let bytes = s.as_bytes();


    let mut result: Vec<&str> = Vec::new();

    loop {

        while current_end < str_length && bytes[current_end] as char == character {
            current_end  += 1
        }
        if current_end == str_length {
            break
        }
        current_start = current_end;
        current_end += 1;

        while current_end < str_length && bytes[current_end] as char  != character {
            current_end  += 1
        }
        // no split string
        if current_start == 0 && current_end == str_length {
            result.push(s);
            break
        }

        let ss: &str = &s[current_start..current_end];
        result.push(ss);
    }
    return result;
}

/// A Python module implemented in Rust.
#[pymodule]
fn strings(py: Python, m: &PyModule) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(split, m)?)?;

    Ok(())
}
```

the caveats should make beating cpython trivial as it is much more flexible, i.e multiple splitting chars and full
utf-8 character, however we are at a slight disadvantage. PyO3 converts the rust &str vector to a list of strings after
the split string vector has been created. Essentially preforming a O(n) operation to create the final python object
behind the scenes. The [cPython version](https://github.com/python/cpython/blob/master/Objects/stringlib/split.h)
which my algorithm is a simplified version of builds the list of python objects directly saving the heap allocations
for the vector.

### The results...

```
length rust:    0.0001556873321533203 sec
length cPython: 0.0000987052917480468 sec
```
> This is for 681 words separated by spaces and containing only US-ASCII characters on a Octa-Core, 8 x 2,4 GHz Turbo with 16gb Ram

Close but no cigar, this is the best I can do with the hour I spent building the module. There likely is a be better way to
get string length than from the char iterator created in line 10 but as this function runs once it is definitely not the limiting factor
to matching the speed of cPython implementation. Nonetheless, given the benefits of the rust development experience, rust can
become a serious contender for python extensions.

