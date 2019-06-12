package slidingwindow

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"testing"
	"time"
)

var test int

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func shouldExecute(tn string) bool {
	exclude := []string{}
	execute := []string{}

	if len(exclude) > 0 {
		return !contains(exclude, tn)
	}

	if len(execute) > 0 {
		return contains(execute, tn)
	}

	return true
}

func precondition(t int) (bool, error) {
	var tn string
	ptr, _, _, ok := runtime.Caller(1)
	if !ok {
		return false, fmt.Errorf("Test %d failed to check precondition", t)
	}

	tn = runtime.FuncForPC(ptr).Name()
	i := strings.LastIndex(tn, ".")
	if i >= 0 {
		tn = tn[i+1:]
	}

	log.Printf("Test %s (%d)", tn, t)

	if !shouldExecute(tn) {
		log.Println("skipped")
		return false, nil
	}

	return true, nil
}

func TestLimit(t *testing.T) {
	test = 1
	testcase := 1

	if b, e := precondition(test); e != nil {
		log.Println(e)
		t.FailNow()
	} else if !b {
		t.SkipNow()
	}

	limit, e := NewLimit(10, 10)
	if e != nil {
		t.Errorf("should be able to create limit (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	if _, e := NewLimit(10, 0); e == nil {
		t.Errorf("should be not able to create limit with '0' windowsize (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	if limit.Check(10) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++

	time.Sleep(6 * time.Second)

	if limit.Check(10) {
		t.Errorf("should not have been allowed (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	time.Sleep(6 * time.Second)

	if limit.Check(10) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
}

func TestNamedLimit(t *testing.T) {

	test = 2
	testcase := 1

	if b, e := precondition(test); e != nil {
		log.Println(e)
		t.FailNow()
	} else if !b {
		t.SkipNow()
	}

	if e := NewNamedLimit("limit1", 10, 10); e != nil {
		t.Errorf("should be able to create 'limit1' (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	if e := NewNamedLimit("limit1", 10, 10); e != nil {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not be able to create 'limit1' (%d.%d)", test, testcase)
	}
	testcase++

	if b, e := Check("limit1", 10); e != nil {
		t.Errorf("'limit1' should be available (%d.%d)", test, testcase)
	} else if b {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++

	if _, e := Check("limit2", 10); e != nil {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("'limit2' should not be available (%d.%d)", test, testcase)
	}
	testcase++

	if b, e := Check("limit1", 10); e != nil {
		t.Errorf("'limit1' should be available (%d.%d)", test, testcase)
	} else if b {
		t.Errorf("should not have been allowed (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	if e := Reset("limit1"); e != nil {
		t.Errorf("'limit1' should be available (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	if b, e := Check("limit1", 10); e != nil {
		t.Errorf("'limit' should be available (%d.%d)", test, testcase)
	} else if b {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++

}

func TestUnlimited(t *testing.T) {
	test = 3
	testcase := 1

	if b, e := precondition(test); e != nil {
		log.Println(e)
		t.FailNow()
	} else if !b {
		t.SkipNow()
	}

	limit, e := NewLimit(0, 0)
	if e != nil {
		log.Println(e)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	if limit.Check(10000) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++

	if limit.Check(10000) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}

}

func TestTiming(t *testing.T) {
	test = 4
	testcase := 1

	if b, e := precondition(test); e != nil {
		log.Println(e)
		t.FailNow()
	} else if !b {
		t.SkipNow()
	}

	limit, e := NewLimit(10, 10)
	if e != nil {
		t.Errorf("should be able to create limit (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	if limit.Check(8) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++

	if limit.Check(1) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++

	if limit.Check(1) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++

	if limit.Check(1) {
		t.Errorf("should not been allowed (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	time.Sleep(1 * time.Second)
	if limit.Check(1) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++

	if limit.Check(1) {
		t.Errorf("should not been allowed (%d.%d)", test, testcase)
	} else {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	}
	testcase++

	time.Sleep(3 * time.Second)
	if limit.Check(3) {
		log.Printf("testcase %d.%d succeeded", test, testcase)
	} else {
		t.Errorf("should not have failed (%d.%d)", test, testcase)
	}
	testcase++
}
