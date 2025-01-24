# goink

An implementation of the ink runtime in go

## Undocumented Runtime gotchas

1. bools are a type that can be returned, not just 1 an 0 in the form of ints and floats

2. '<>' operator is not documented in the runtime, it should remove excessive line breaks between strings
