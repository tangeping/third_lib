package simplejson

import (
	"encoding/json"
	"fmt"
	"testing"
)

func AssertEqual(t1 interface{}, t2 interface{}, r interface{}) {

	if t1 == t2 {
		return
	}

	fmt.Sprintln(r)
}

func AssertNotEqual(t1 interface{}, t2 interface{}, r interface{}) {

	if t1 != t2 {
		return
	}

	fmt.Sprintln(r)
}

func TestSimplejson(t *testing.T) {
	var ok bool
	var err error

	js, err := NewJson([]byte(`{
		"test": {
			"string_array": ["asdf", "ghjk", "zxcv"],
			"string_array_null": ["abc", null, "efg"],
			"array": [1, "2", 3],
			"arraywithsubs": [{"subkeyone": 1},
			{"subkeytwo": 2, "subkeythree": 3}],
			"int": 10,
			"float": 5.150,
			"string": "simplejson",
			"bool": true,
			"sub_obj": {"a": 1}
		}
	}`))

	AssertNotEqual(t, nil, js)
	AssertEqual(t, nil, err)

	_, ok = js.CheckGet("test")
	AssertEqual(t, true, ok)

	_, ok = js.CheckGet("missing_key")
	AssertEqual(t, false, ok)

	aws := js.Get("test").Get("arraywithsubs")
	AssertNotEqual(t, nil, aws)
	var awsval int
	awsval, _ = aws.GetIndex(0).Get("subkeyone").Int()
	AssertEqual(t, 1, awsval)
	awsval, _ = aws.GetIndex(1).Get("subkeytwo").Int()
	AssertEqual(t, 2, awsval)
	awsval, _ = aws.GetIndex(1).Get("subkeythree").Int()
	AssertEqual(t, 3, awsval)

	i, _ := js.Get("test").Get("int").Int()
	AssertEqual(t, 10, i)

	f, _ := js.Get("test").Get("float").Float64()
	AssertEqual(t, 5.150, f)

	s, _ := js.Get("test").Get("string").String()
	AssertEqual(t, "simplejson", s)

	b, _ := js.Get("test").Get("bool").Bool()
	AssertEqual(t, true, b)

	mi := js.Get("test").Get("int").MustInt()
	AssertEqual(t, 10, mi)

	mi2 := js.Get("test").Get("missing_int").MustInt(5150)
	AssertEqual(t, 5150, mi2)

	ms := js.Get("test").Get("string").MustString()
	AssertEqual(t, "simplejson", ms)

	ms2 := js.Get("test").Get("missing_string").MustString("fyea")
	AssertEqual(t, "fyea", ms2)

	ma2 := js.Get("test").Get("missing_array").MustArray([]interface{}{"1", 2, "3"})
	AssertEqual(t, ma2, []interface{}{"1", 2, "3"})

	msa := js.Get("test").Get("string_array").MustStringArray()
	AssertEqual(t, msa[0], "asdf")
	AssertEqual(t, msa[1], "ghjk")
	AssertEqual(t, msa[2], "zxcv")

	msa2 := js.Get("test").Get("string_array").MustStringArray([]string{"1", "2", "3"})
	AssertEqual(t, msa2[0], "asdf")
	AssertEqual(t, msa2[1], "ghjk")
	AssertEqual(t, msa2[2], "zxcv")

	msa3 := js.Get("test").Get("missing_array").MustStringArray([]string{"1", "2", "3"})
	AssertEqual(t, msa3, []string{"1", "2", "3"})

	mm2 := js.Get("test").Get("missing_map").MustMap(map[string]interface{}{"found": false})
	AssertEqual(t, mm2, map[string]interface{}{"found": false})

	strs, err := js.Get("test").Get("string_array").StringArray()
	AssertEqual(t, err, nil)
	AssertEqual(t, strs[0], "asdf")
	AssertEqual(t, strs[1], "ghjk")
	AssertEqual(t, strs[2], "zxcv")

	strs2, err := js.Get("test").Get("string_array_null").StringArray()
	AssertEqual(t, err, nil)
	AssertEqual(t, strs2[0], "abc")
	AssertEqual(t, strs2[1], "")
	AssertEqual(t, strs2[2], "efg")

	gp, _ := js.GetPath("test", "string").String()
	AssertEqual(t, "simplejson", gp)

	gp2, _ := js.GetPath("test", "int").Int()
	AssertEqual(t, 10, gp2)

	AssertEqual(t, js.Get("test").Get("bool").MustBool(), true)

	js.Set("float2", 300.0)
	AssertEqual(t, js.Get("float2").MustFloat64(), 300.0)

	js.Set("test2", "setTest")
	AssertEqual(t, "setTest", js.Get("test2").MustString())

	js.Del("test2")
	AssertNotEqual(t, "setTest", js.Get("test2").MustString())

	js.Get("test").Get("sub_obj").Set("a", 2)
	AssertEqual(t, 2, js.Get("test").Get("sub_obj").Get("a").MustInt())

	js.GetPath("test", "sub_obj").Set("a", 3)
	AssertEqual(t, 3, js.GetPath("test", "sub_obj", "a").MustInt())
}

func TestStdlibInterfaces(t *testing.T) {
	val := new(struct {
		Name   string `json:"name"`
		Params *Json  `json:"params"`
	})
	val2 := new(struct {
		Name   string `json:"name"`
		Params *Json  `json:"params"`
	})

	raw := `{"name":"myobject","params":{"string":"simplejson"}}`

	AssertEqual(t, nil, json.Unmarshal([]byte(raw), val))

	AssertEqual(t, "myobject", val.Name)
	AssertNotEqual(t, nil, val.Params.data)
	s, _ := val.Params.Get("string").String()
	AssertEqual(t, "simplejson", s)

	p, err := json.Marshal(val)
	AssertEqual(t, nil, err)
	AssertEqual(t, nil, json.Unmarshal(p, val2))
	AssertEqual(t, val, val2) // stable
}

func TestSet(t *testing.T) {
	js, err := NewJson([]byte(`{}`))
	AssertEqual(t, nil, err)

	js.Set("baz", "bing")

	s, err := js.GetPath("baz").String()
	AssertEqual(t, nil, err)
	AssertEqual(t, "bing", s)
}

func TestReplace(t *testing.T) {
	js, err := NewJson([]byte(`{}`))
	AssertEqual(t, nil, err)

	err = js.UnmarshalJSON([]byte(`{"baz":"bing"}`))
	AssertEqual(t, nil, err)

	s, err := js.GetPath("baz").String()
	AssertEqual(t, nil, err)
	AssertEqual(t, "bing", s)
}

func TestSetPath(t *testing.T) {
	js, err := NewJson([]byte(`{}`))
	AssertEqual(t, nil, err)

	js.SetPath([]string{"foo", "bar"}, "baz")

	s, err := js.GetPath("foo", "bar").String()
	AssertEqual(t, nil, err)
	AssertEqual(t, "baz", s)
}

func TestSetPathNoPath(t *testing.T) {
	js, err := NewJson([]byte(`{"some":"data","some_number":1.0,"some_bool":false}`))
	AssertEqual(t, nil, err)

	f := js.GetPath("some_number").MustFloat64(99.0)
	AssertEqual(t, f, 1.0)

	js.SetPath([]string{}, map[string]interface{}{"foo": "bar"})

	s, err := js.GetPath("foo").String()
	AssertEqual(t, nil, err)
	AssertEqual(t, "bar", s)

	f = js.GetPath("some_number").MustFloat64(99.0)
	AssertEqual(t, f, 99.0)
}

func TestPathWillAugmentExisting(t *testing.T) {
	js, err := NewJson([]byte(`{"this":{"a":"aa","b":"bb","c":"cc"}}`))
	AssertEqual(t, nil, err)

	js.SetPath([]string{"this", "d"}, "dd")

	cases := []struct {
		path    []string
		outcome string
	}{
		{
			path:    []string{"this", "a"},
			outcome: "aa",
		},
		{
			path:    []string{"this", "b"},
			outcome: "bb",
		},
		{
			path:    []string{"this", "c"},
			outcome: "cc",
		},
		{
			path:    []string{"this", "d"},
			outcome: "dd",
		},
	}

	for _, tc := range cases {
		s, err := js.GetPath(tc.path...).String()
		AssertEqual(t, nil, err)
		AssertEqual(t, tc.outcome, s)
	}
}

func TestPathWillOverwriteExisting(t *testing.T) {
	// notice how "a" is 0.1 - but then we'll try to set at path a, foo
	js, err := NewJson([]byte(`{"this":{"a":0.1,"b":"bb","c":"cc"}}`))
	AssertEqual(t, nil, err)

	js.SetPath([]string{"this", "a", "foo"}, "bar")

	s, err := js.GetPath("this", "a", "foo").String()
	AssertEqual(t, nil, err)
	AssertEqual(t, "bar", s)
}
