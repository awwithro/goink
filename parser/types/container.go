package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Container struct {
	ParentContainer *Container
	Contents        []Acceptor
	Name            string
	Flag            byte
	SubContainers   map[string]*Container
}

func NewContainer(name string, parent *Container) *Container {
	c := &Container{
		Name:            name,
		Contents:        make([]Acceptor, 0),
		SubContainers:   make(map[string]*Container),
		ParentContainer: parent,
	}
	return c
}

func (c *Container) Accept(v Visitor) {
	v.VisitContainer(c)
}

func (c *Container) UnmarshalJSON(p []byte) error {
	var raw []any
	if err := json.Unmarshal(p, &raw); err != nil {
		return err
	}
	err := c.unmarshalContainer(raw)
	return err
}

func (c *Container) unmarshalString(str string) error {
	if strings.HasPrefix(str, "^") || str == "\n" {
		c.Contents = append(c.Contents,
			StringVal(strings.TrimPrefix(str, "^")))
	} else if cmd, ok := IsControlCommand(str); ok {
		c.Contents = append(c.Contents, cmd)
	} else if op, ok := IsOperator(str); ok {
		c.Contents = append(c.Contents, Operator(op))
	}
	return nil
}

func (c *Container) unmarshalContainer(cnt []any) error {
	// Final element is null or a specific map
	if finalElement, ok := cnt[len(cnt)-1].(map[string]any); ok {
		c.parseFinalElement(finalElement)
	}
	// Skipping the last element since we just parsed it
	for x := 0; x < len(cnt)-1; x++ {
		val := cnt[x]
		switch typ := val.(type) {
		case []any:
			subContainer := NewContainer("", c)
			subContainer.unmarshalContainer(typ)
			c.Contents = append(c.Contents, subContainer)
		case string:
			c.unmarshalString(typ)
		// All numbers get converted to floats. strconv dance
		// to figure out if we can use an int
		case float64:
			str := strconv.FormatFloat(typ, 'f', -1, 64)
			if num, err := strconv.Atoi(str); err == nil {
				c.Contents = append(c.Contents, IntVal(num))
			} else {
				c.Contents = append(c.Contents, FloatVal(typ))
			}
		case map[string]any:
			c.unmarshalMaps(typ)
		case bool:
			c.Contents = append(c.Contents, BoolVal(typ))
		default:
			log.Error("Unrecognized Container element ", val)
		}
	}
	return nil
}

func (c *Container) parseFinalElement(obj map[string]any) {
	for k, v := range obj {
		switch k {
		case "#n":
			name := v.(string)
			c.Name = name
		case "#f":
			flag := v.(float64)
			c.Flag = byte(flag)
		default:
			if cnt, ok := v.([]any); ok {
				subContainer := NewContainer(k, c)
				subContainer.unmarshalContainer(cnt)
				c.SubContainers[k] = subContainer
			} else {
				log.Panic("Unrecognized Final Element ", k, v)
			}
		}
	}
}

func (c *Container) unmarshalMaps(obj map[string]any) {
	for k, v := range obj {
		switch k {
		case "^->":
			target := v.(string)
			c.Contents = append(c.Contents, DivertTarget(target))
		case "^var":
			target := v.(string)
			ptr := NewVariablePointer(target)
			if i, ok := obj["ci"]; ok {
				idx := i.(float64)
				ptr.ContextIndex = int(idx)
			}
			c.Contents = append(c.Contents, ptr)
		case "->":
			target := v.(string)
			var conditional bool
			var variable bool
			if i, ok := obj["c"]; ok {
				conditional = i.(bool)
			}
			if i, ok := obj["var"]; ok {
				variable = i.(bool)
			}
			if variable {
				c.Contents = append(c.Contents, VariableDivert{
					Name:        target,
					Conditional: conditional,
				})
			} else {
				c.Contents = append(c.Contents, Divert{
					Path:        Path(target),
					Conditional: conditional,
				})
			}

		case "f()":
			target := v.(string)
			var conditional bool
			if i, ok := obj["c"]; ok {
				conditional = i.(bool)
			}
			c.Contents = append(c.Contents, FunctionDivert{
				Path:        Path(target),
				Conditional: conditional,
			})
		case "->t->":
			target := v.(string)
			var conditional bool
			if i, ok := obj["c"]; ok {
				conditional = i.(bool)
			}
			c.Contents = append(c.Contents, TunnelDivert{
				Path:        Path(target),
				Conditional: conditional,
			})
		case "x()":
			target := v.(string)
			var conditional bool
			var args int
			if i, ok := obj["c"]; ok {
				conditional = i.(bool)
			}
			if i, ok := obj["exArgs"]; ok {
				conditional = i.(bool)
			}
			c.Contents = append(c.Contents, ExternalFunctionDivert{
				Path:        Path(target),
				Conditional: conditional,
				Args:        args,
			})
		case "*":
			path := v.(string)
			var flag byte
			if f, ok := obj["flg"]; ok {
				flag = byte(f.(float64))
			}
			c.Contents = append(c.Contents, ChoicePoint{
				Path: Path(path),
				Flag: flag,
			})
		case "VAR?":
			val := v.(string)
			c.Contents = append(c.Contents, VarRef(val))
		case "temp=":
			name := v.(string)
			var re bool
			if r, ok := obj["re"]; ok {
				re = r.(bool)
			}
			c.Contents = append(c.Contents, TempVar{
				Name:     name,
				ReAssign: re,
			})
		case "VAR=":
			name := v.(string)
			var re bool
			if r, ok := obj["re"]; ok {
				re = r.(bool)
			}
			c.Contents = append(c.Contents, GlobalVar{
				Name:     name,
				ReAssign: re,
			})
		case "CNT?":
			val := v.(string)
			c.Contents = append(c.Contents, ReadCount(val))
		// parsed as part of other key
		case "var", "c", "exArgs", "ci", "flg", "re":
			continue
		default:
			log.Warn("Unrecognized key ", k)
		}
	}
}

func (c *Container) RecordVisits() bool {
	return c.Flag&0x1 == 1
}
func (c *Container) RecordTurns() bool {
	return c.Flag&0x2 == 2
}
func (c *Container) CountStartOnly() bool {
	return c.Flag&0x4 == 4
}

func (c *Container) GetNamedContainer(name string) (*Container, error) {
	if cnt, ok := c.SubContainers[name]; ok {
		return cnt, nil
	} else {
		for _, obj := range c.Contents {
			cnt, ok := obj.(*Container)
			if ok && cnt.Name == name {
				return cnt, nil
			}
		}
	}
	return nil, NoNamedContainer(fmt.Errorf("no container named %s found", name))
}

func (c *Container) GetRoot() *Container {
	cnt := c
	for {
		if cnt.ParentContainer == nil {
			return cnt
		}
		cnt = cnt.ParentContainer
	}
}

func (c *Container) PositionInParent() (int, error) {
	if c.ParentContainer != nil {
		_, ok := c.ParentContainer.SubContainers[c.Name]
		if ok {
			return -1, EndOfSubContainer{}
		}
		for x := 0; x < len(c.ParentContainer.Contents); x++ {
			obj := c.ParentContainer.Contents[x]
			if cnt, ok := obj.(*Container); ok {
				if cnt == c {
					return x, nil
				}
			}
		}
	}
	return -1, fmt.Errorf("no Parent")
}

type NoNamedContainer error

type EndOfSubContainer struct{}

func (e EndOfSubContainer) Error() string {
	return "reached end of sub-container"
}
