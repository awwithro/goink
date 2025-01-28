# goink

An implementation of the ink runtime in go

## Undocumented Runtime gotchas

1. bools are a type that can be returned, not just 1 an 0 in the form of ints and floats
2. '<>' operator is not documented in the runtime, it should remove excessive line breaks between strings
3. Global vars get assigned as a special sub container in the root container named "global decl"
4. Visit should be previous visits. First visit to a container should have "0"
5. LIST_MIN, LIST_MAX,LIST_ALL,LIST_COUNT,LIST_VALUE, LIST_INSERT undocumented. All list functions. This also explains the listDefs object. listDefs are map where a key is a name of a list, each item in the list is a kv pair of the value and index.

6. More undocumented operations: POW, FLOOR, CEILING, INT, FLOAT, ?, !?. ^
```
{
    "list": {},
    "origins": [
        "BedKnowledge"
    ]
},
{
    "VAR=": "BedKnowledge
}
```
references a list, fils it with the defined BedKnowledge listDef and assigns it to the globalVar BedKnowledge