package format

import "testing"

func TestColor(t *testing.T) {
	res := Color(Format([]byte(`
{"hello":"world","what":123,
"arr":["1","2",1,2,true,false,null],
"obj":{"key1":null,"ar`+"\x1B[36m"+`Cyanr2":[1,2,3,"123","456"]}}
	`)), nil)
	if string(res) != `{
  [94m"hello"[0m: [92m"world"[0m,
  [94m"what"[0m: [93m123[0m,
  [94m"arr"[0m: [[92m"1"[0m, [92m"2"[0m, [93m1[0m, [93m2[0m, [96mtrue[0m, [96mfalse[0m, [91mnull[0m],
  [94m"obj"[0m: {
    [94m"key1"[0m: [91mnull[0m,
    [94m"ar\u001b[36mCyanr2"[0m: [[93m1[0m, [93m2[0m, [93m3[0m, [92m"123"[0m, [92m"456"[0m]
  }
}
` {
		t.Fatal("invalid output")
	}
}
