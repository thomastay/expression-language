# Parity with spec

1. Floor division (DONE)
1. exponentiation (DONE)
1. Arrays (DONE)
1. Make declaring an external function easier (DONE)
1. to string conversion in VM itself
1. Make Boolean a type in the expression engine so we can print it properly (DONE)
1. Split bytecode struct into pieces (do a SOA optimization) (DONE)
1. Move values to a constant table and make Load instructions read from this table (DONE)
1. Add Immediate instructions for cases where one side is a constant

# Features

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
