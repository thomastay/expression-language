# Parity with spec

1. Floor division (DONE)
1. exponentiation (DONE)
1. Arrays
1. to string conversion in VM itself

# Features

1. Split bytecode struct into pieces (do a SOA optimization)
   As part of this, also move values to a constant table and make Load instructions read from this table
1. Constant folding (on AST level)
   Split vm math module into a runtime module so compiler can use it as well
   Fold +, -, \*, %, /, //
   Possibly write a walk() and visit() interface like in the Expr library
   Also fold conditionals if the conditional is a constant
1. Jump chaining (a peephole optimization)
1. Implement local variables
1. Implement SSA based IR
1. Implement CSE and DCE (relies on the above 2)
1. Implement closures
