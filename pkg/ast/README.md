# Syntax diagram

```
Start : E
E : V
  | E op E
  | uop E
  | ( E )
  | E ? E : E
  | E [ E ]
  | Call
  | KArr
  | EFieldAccess

Call: Ident ( EList )
    | E . Ident ( EList )

EFieldAccess : E . Ident

EList : E ',' EList
      | E

V : Ident | Int | Float | KStr | true | false

KArr  : [ VList ]
VList : V ',' VList
      | V
```

## Notes

Expression parsing is with the Pratt parser
op is any infix op
uop is any prefix unary op
Kaitai struct has no postfix ops (phew)
KString means Constant String
