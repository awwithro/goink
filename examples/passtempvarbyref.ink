->test

== test ==
~ temp x = 1
initial x {x}
~baz(x)
-> DONE


== function baz(ref var)==
~var +=1
bar var {var}
~ barref(var)
~return

== function barref(ref var) ==
~var +=1
barref var {var}
~return