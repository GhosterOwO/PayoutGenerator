package job

import (
	"regexp"
	"testing"
)

func TestSomething(t *testing.T) {
	s := `
12 7 8 
1 4 18 
12 17 15 
19 14 1 
16 19 14 
4 2 1 
8 4 7 
8 3 17 
17 20 6 
1 2 17 
9 5 12 
12 15 1 
1 9 19 
14 1 18 
11 18 9 
4 7 13 
11 4 17 
19 6 11 
3 4 17 
1 7 9 
1 3 15 
3 10 19 
16 14 4 
8 16 11 
20 14 12 
9 1 15 
14 1 8 
18 13 14 
7 5 10 
14 0 3 
20 9 0 
15 8 4 
10 9 15 
10 4 11 
	`
	pattern := "[[:space:]]0[[:space:]]"
	re := regexp.MustCompile("(?m)^.*" + pattern + ".*$[\r\n]+")
	chunk := re.ReplaceAllString(s, "")
	if match, _ := regexp.Match(pattern, []byte(chunk)); match {
		t.Error(chunk)
	}
}
