package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/awwithro/goink/pkg/parser/types"
)

func TestContainer(t *testing.T) {
	rawJson := `{"inkVersion":21,"root":[["^Once upon a time...","\n",["ev",{"^->":"0.2.$r1"},{"temp=":"$r"},"str",{"->":".^.s"},[{"#n":"$r1"}],"/str","/ev",{"*":"0.c-0","flg":18},{"s":["^There were two choices.",{"->":"$r","var":true},null]}],["ev",{"^->":"0.3.$r1"},{"temp=":"$r"},"str",{"->":".^.s"},[{"#n":"$r1"}],"/str","/ev",{"*":"0.c-1","flg":18},{"s":["^There were four lines of content.",{"->":"$r","var":true},null]}],{"c-0":["ev",{"^->":"0.c-0.$r2"},"/ev",{"temp=":"$r"},{"->":"0.2.s"},[{"#n":"$r2"}],"\n",{"->":"0.g-0"},{"#f":5}],"c-1":["ev",{"^->":"0.c-1.$r2"},"/ev",{"temp=":"$r"},{"->":"0.3.s"},[{"#n":"$r2"}],"\n",{"->":"0.g-0"},{"#f":5}],"g-0":["^They lived happily ever after.","\n","end",["done",{"#f":5,"#n":"g-1"}],{"#f":5}]}],"done",{"#f":1}],"listDefs":{}}`
	i := Parse([]byte(rawJson))
	printContainer(&i.Root, 0)
}

func printContainer(cnt *types.Container, depth int) {
	for _, val := range cnt.Contents {
		fmt.Print(strings.Repeat("-", depth))
		switch typ := val.(type) {
		case *types.Container:
			fmt.Printf(" Container \"%s\" %v\n", typ.Name, typ.Flag)
			printContainer(typ, depth+1)
		case types.ControlCommand:
			fmt.Printf(" ControlCommand: %v\n", typ)
		case types.StringVal:
			fmt.Printf(" %s\n", strings.ReplaceAll(typ.String(), "\n", "\\n"))
		case types.Divert:
			fmt.Printf(" Divert -> %s\n", typ.Path)
		case types.VariableDivert:
			fmt.Printf(" Divert -> %s\n", typ.Name)
		case types.DivertTarget:
			fmt.Printf(" DivertTarget: %s\n", typ)
		default:
			fmt.Printf(" %s\n", reflect.TypeOf(typ))
		}
	}
	if len(cnt.SubContainers) > 0 {
		fmt.Print(strings.Repeat("-", depth))
		fmt.Println(" SubContainers")
		for _, v := range cnt.SubContainers {
			printContainer(v, depth+1)
		}
	}
}
