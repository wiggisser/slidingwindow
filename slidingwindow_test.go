package slidingwindow

import (
	"fmt"
	"testing"
	"time"
)

var testcase int

func TestLimit(t *testing.T) {
	testcase = 1
	test := 1
	fmt.Printf("Testcase TestLimit (%d)\r\n**********************\r\n", testcase)

	limit, e := NewLimit(10, 10)
	if e != nil {
		fmt.Println(e)
		t.Errorf("should be able to create limit (%d.%d)", testcase, test)
	}
	test++

	if _, e := NewLimit(10, 0); e == nil {
		t.Errorf("should be not able to create limit with '0' windowsize (%d.%d)", testcase, test)
	}
	test++

	if limit.Check(10) {
		fmt.Println("allowed")
	} else {
		fmt.Println("quota exceeded")
		t.Errorf("Should not have failed (%d.%d)", testcase, test)
	}
	test++

	time.Sleep(6 * time.Second)

	if limit.Check(10) {
		fmt.Println("allowed")
		t.Errorf("Should not have been allowed (%d.%d)", testcase, test)
	} else {
		fmt.Println("quota exceeded")
	}
	test++

	time.Sleep(6 * time.Second)

	if limit.Check(10) {
		fmt.Println("allowed")
	} else {
		fmt.Println("quota exceeded")
		t.Errorf("Should not have failed (%d.%d)", testcase, test)
	}
}

func TestNamedLimit(t *testing.T) {

	testcase = 2
	test := 1

	fmt.Printf("Testcase TestNamedLimit (%d)\r\n***************************\r\n", testcase)

	if e := NewNamedLimit("limit1", 10, 10); e != nil {
		t.Errorf("should be able to create 'limit1' (%d.%d)", testcase, test)
		fmt.Println(e)
	} else {
		fmt.Println("named limit 'limit1' created")
	}
	test++

	if e := NewNamedLimit("limit1", 10, 10); e != nil {
		fmt.Println(e)
	} else {
		t.Errorf("should not be able to create 'limit1' (%d.%d)", testcase, test)
		fmt.Println("named limit 'limit1' created")
	}
	test++

	if b, e := Check("limit1", 10); e != nil {
		t.Errorf("'limit1' should be available (%d.%d)", testcase, test)
		fmt.Println(e)
	} else if b {
		fmt.Println("allowed")
	} else {
		t.Errorf("should not have failed (%d.%d)", testcase, test)
		fmt.Println("quota exceeded")
	}
	test++

	if b, e := Check("limit2", 10); e != nil {
		fmt.Println(e)
	} else if b {
		t.Errorf("'limit2' should not be available (%d.%d)", testcase, test)
		fmt.Println("allowed")
	} else {
		t.Errorf("'limit2' should not be available (%d.%d)", testcase, test)
		fmt.Println("quota exceeded")
	}
	test++

	if b, e := Check("limit1", 10); e != nil {
		t.Errorf("'limit1' should be available (%d.%d)", testcase, test)
		fmt.Println(e)
	} else if b {
		t.Errorf("Should not have been allowed (%d.%d)", testcase, test)
		fmt.Println("allowed")
	} else {
		fmt.Println("quota exceeded")
	}
	test++

	if e := Reset("limit1"); e != nil {
		t.Errorf("'limit1' should be available (%d.%d)", testcase, test)
		fmt.Println(e)
	} else {
		fmt.Println("reset 'limit1' succeeded")
	}
	test++

	if b, e := Check("limit1", 10); e != nil {
		t.Errorf("'limit' should be available (%d.%d)", testcase, test)
		fmt.Println(e)
	} else if b {
		fmt.Println("allowed")
	} else {
		t.Errorf("should not have failed (%d.%d)", testcase, test)
		fmt.Println("quota exceeded")
	}
	test++

}

func TestUnlimited(t *testing.T) {
	testcase = 3
	test := 1

	fmt.Printf("Testcase TestUnlimited (%d)\r\n**************************\r\n", testcase)

	limit, e := NewLimit(0, 0)
	if e != nil {
		fmt.Println(e)
		t.Errorf("should be able to create unlimited limit (%d.%d)", testcase, test)
	}
	test++

	if limit.Check(10000) {
		fmt.Println("allowed")
	} else {
		t.Errorf("should not have failed (%d.%d)", testcase, test)
	}
	test++

	if limit.Check(10000) {
		fmt.Println("allowed")
	} else {
		t.Errorf("should not have failed (%d.%d)", testcase, test)
	}

}
