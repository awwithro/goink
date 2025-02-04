LIST test = one, two, (three), (six)
LIST other = four = 4, five = 5, three = 3
LIST third = one=2, seven,eight
{test}
~test = test.one
{test == test.one }
get the representation of a list object: {test.one}
VAR idx = test.two
{idx}
get the value of a list element: {LIST_VALUE(test.three)}
~test=(third.one,two,four)
compare two list objects: {test.three == other.three}
~other=other.three
{test(1)}
Pre-Increment {other}
~other ++
Post increment {other}
{LIST_MIN(other)}