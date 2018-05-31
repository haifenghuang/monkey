if exists("b:current_syntax") 
  finish 
endif 


syn case match 

syn keyword     monkeyDirective         include
syn keyword     monkeyDeclaration       let


hi def link     monkeyDirective         Statement
hi def link     monkeyDeclaration       Type 


syn keyword     monkeyStatement         return let spawn defer struct enum using
syn keyword     monkeyException         try catch finally throw
syn keyword     monkeyConditional       if else elseif elsif elif unless where and or case in is
syn keyword     monkeyRepeat            do while for break continue grep map
syn keyword     monkeyBranch            break continue
syn keyword     monkeyClass             class new property get set default this parent static public private protected interface

hi def link     monkeyStatement         Statement
hi def link     monkeyClass             Statement
hi def link     monkeyConditional       Conditional
hi def link     monkeyBranch            Conditional
hi def link     monkeyLabel             Label
hi def link     monkeyRepeat            Repeat


syn match       monkeyDeclaration       /\<fn\>/
syn match       monkeyDeclaration       /^fn\>/

syn match comment "#.*$"    contains=@Spell,monkeyCommentTodo
syn match comment "\/\/.*$" contains=@Spell,monkeyCommentTodo

syn keyword     monkeyCast              int str float array


hi def link     monkeyCast              Type


syn keyword     monkeyBuiltins          len 
syn keyword     monkeyBuiltins          println print stdin stdout stderr
syn keyword     monkeyBoolean           true false
syn keyword     monkeyNull              nil

hi def link     monkeyBuiltins          Keyword 
hi def link     monkeyNull              Keyword
hi def link     monkeyBoolean           Boolean


" Comments; their contents 
syn keyword     monkeyTodo              contained TODO FIXME XXX BUG 
syn cluster     monkeyCommentGroup      contains=monkeyTodo 
syn region      monkeyComment           start="#" end="$" contains=@monkeyCommentGroup,@Spell,@monkeyTodo


hi def link     monkeyComment           Comment 
hi def link     monkeyTodo              Todo 


" monkey escapes 
syn match       monkeyEscapeOctal       display contained "\\[0-7]\{3}" 
syn match       monkeyEscapeC           display contained +\\[abfnrtv\\'"]+ 
syn match       monkeyEscapeX           display contained "\\x\x\{2}" 
syn match       monkeyEscapeU           display contained "\\u\x\{4}" 
syn match       monkeyEscapeBigU        display contained "\\U\x\{8}" 
syn match       monkeyEscapeError       display contained +\\[^0-7xuUabfnrtv\\'"]+ 


hi def link     monkeyEscapeOctal       monkeySpecialString 
hi def link     monkeyEscapeC           monkeySpecialString 
hi def link     monkeyEscapeX           monkeySpecialString 
hi def link     monkeyEscapeU           monkeySpecialString 
hi def link     monkeyEscapeBigU        monkeySpecialString 
hi def link     monkeySpecialString     Special 
hi def link     monkeyEscapeError       Error 
hi def link     monkeyException		Exception

" Strings and their contents 
syn cluster     monkeyStringGroup       contains=monkeyEscapeOctal,monkeyEscapeC,monkeyEscapeX,monkeyEscapeU,monkeyEscapeBigU,monkeyEscapeError 
syn region      monkeyString            start=+"+ skip=+\\\\\|\\"+ end=+"+ contains=@monkeyStringGroup
syn region      monkeyRegExString       start=+/[^/*]+me=e-1 skip=+\\\\\|\\/+ end=+/\s*$+ end=+/\s*[;.,)\]}]+me=e-1 oneline
syn region      monkeyRawString         start=+`+ end=+`+ 


hi def link     monkeyString            String 
hi def link     monkeyRawString         String 
hi def link     monkeyRegExString       String

" Characters; their contents 
syn cluster     monkeyCharacterGroup    contains=monkeyEscapeOctal,monkeyEscapeC,monkeyEscapeX,monkeyEscapeU,monkeyEscapeBigU 
syn region      monkeyCharacter         start=+'+ skip=+\\\\\|\\'+ end=+'+ contains=@monkeyCharacterGroup 


hi def link     monkeyCharacter         Character 


" Regions 
syn region      monkeyBlock             start="{" end="}" transparent fold 
syn region      monkeyParen             start='(' end=')' transparent 


" Integers 
syn match       monkeyDecimalInt        "\<\d\+\([Ee]\d\+\)\?\>" 
syn match       monkeyHexadecimalInt    "\<0x\x\+\>" 
syn match       monkeyOctalInt          "\<0\o\+\>" 
syn match       monkeyOctalError        "\<0\o*[89]\d*\>" 


hi def link     monkeyDecimalInt        Integer
hi def link     monkeyHexadecimalInt    Integer
hi def link     monkeyOctalInt          Integer
hi def link     Integer                 Number

" Floating point 
syn match       monkeyFloat             "\<\d\+\.\d*\([Ee][-+]\d\+\)\?\>" 
syn match       monkeyFloat             "\<\.\d\+\([Ee][-+]\d\+\)\?\>" 
syn match       monkeyFloat             "\<\d\+[Ee][-+]\d\+\>" 


hi def link     monkeyFloat             Float 
"hi def link     monkeyImaginary         Number 


if exists("monkey_fold")
    syn match	monkeyFunction	"\<fn\>"
    syn region	monkeyFunctionFold	start="\<fn\>.*[^};]$" end="^\z1}.*$" transparent fold keepend

    syn sync match monkeySync	grouphere monkeyFunctionFold "\<fn\>"
    syn sync match monkeySync	grouphere NONE "^}"

    setlocal foldmethod=syntax
    setlocal foldtext=getline(v:foldstart)
else
    syn keyword monkeyFunction	fn
    syn match	monkeyBraces	"[{}\[\]]"
    syn match	monkeyParens	"[()]"
endif

syn sync fromstart
syn sync maxlines=100

hi def link monkeyFunction		Function
hi def link monkeyBraces		Function

let b:current_syntax = "monkey" 
