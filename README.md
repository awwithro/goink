# goink

An implementation of the ink runtime in go

## Undocumented Runtime gotchas

1. bools are a type that can be returned, not just 1 an 0 in the form of ints and floats
2. '<>' operator is not documented in the runtime, it should remove excessive line breaks between strings
3. Global vars get assigned as a special sub container in the root container named "global decl"
4. Visit should be previous visits. First visit to a container should have "0"
5. LIST_MIN, LIST_MAX,LIST_ALL,LIST_COUNT,LIST_VALUE, LIST_INSERT undocumented. All list functions. This also explains the listDefs object. listDefs are map where a key is a name of a list, each item in the list is a kv pair of the value and index.

6. More undocumented operations: POW, FLOOR, CEILING, INT, FLOAT, ?, !?. ^
7. listInt pops a list name and an index, prints that name
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

## List Notes
* List and a List Val are two distinct values.
* Need to create a global val for each key
  * if a key is unique across all lists, the global val is the same as the key
  * if a key in not unique, it is created as list.key
* lists themselves are being set as global vars automatically in the 'global decl block
* Equality is based on fqdn. ie foo.bar != baz.bar env if both are set to "1"
* Str is based on key, both foo.bar and baz.bar print "bar"
* Var ref is two the ListVal itself, not the str or map val of the item
* ListInt is odd, it references a list with a string value rather than a proper var ref.
* L^ is a list intersection
* LIST_RANDOM should return a list with the random item not a random item.
* Really all list functions should operate on Lists and not ListItems
* List_Range can operate on all lists, i'm guessing use the int vals in the range func?
* Updates needed to make threads work
    * Need another stack for threads to push/pop from
    * Need to modify DONE handling so we pop a thread if we are in one and then continue evaluation
    * Currently we hit the end of the thread and end evaluation meaning that only one thread will be evaluated at a time
    * Thread contains:
        * TempVars
        * Mode? Could work the same as i have it by saving the mode when returning
        * Pointer to location