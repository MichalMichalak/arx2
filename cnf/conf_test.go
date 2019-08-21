package cnf

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestVal(t *testing.T) {
	type sdef struct {
		id int
	}
	s := sdef{1}
	v := reflect.ValueOf(s)
	vi := v.Interface()
	fmt.Printf("%#v\n", vi)
	process(vi)
}

func process(i interface{}) {
	it := reflect.TypeOf(i)
	fmt.Printf("%#v\n", it)

	iv := reflect.ValueOf(i)
	fmt.Printf("%#v\n", iv)

}

func TestUnm(t *testing.T) {
	type user struct {
		id int
	}
	c := user{}
	err := json.Unmarshal([]byte(`{}`), &c)
	require.NoError(t, err)
	fmt.Println(c)
}

func TestCast(t *testing.T) {
	iii, err := convertToInt("1a")
	fmt.Println(iii, err)
}
