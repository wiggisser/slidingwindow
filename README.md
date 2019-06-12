# Slididing Window

Provides a simple implementation of requestlimits in a sliding time window.

Define a allowed quota per timerange and check every request, wheter it is still allowed in the sliding window. 

## Installation

    go get github.com/wiggisser/slidingwindow

## Examples

### Limit variable

    package main

    import (
        "github.com/wiggisser/slidingwindow"
        "fmt"
        "time"
    }


    func main() {
        //allow 10 requests per 10 seconds
        limit := slidingwindow.NewLimit(10, 10)

        if limit.Check(5) {
            fmt.Println("allowed")
        } else {
            fmt.Println("quota exceeded")
        }

        limit.Reset()

        if limit.Check(5) {
            fmt.Println("allowed")
        } else {
            fmt.Println("quota exceeded")
        }

        time.Sleep(10 * time.Second)

        if limit.Check(5) {
            fmt.Println("allowed")
        } else {
            fmt.Println("quota exceeded")
        }

    }

### Named Limits

    package main

    import (
        "github.com/wiggisser/slidingwindow"
        "fmt"
        "time"
    }


    func main() {
        if e := slidingwindow.NewNamedLimit("limit1", 10, 10); e != nil {
            fmt.Println(e)
        } else {
            fmt.Println("named limit 'limit1' created")
        }

        if b, e := slidingwindow.Check("limit1", 10); e != nil {
            fmt.Println(e);
        } else if b {
            fmt.Println("allowed")
        } else {
            fmt.Println("quota exceeded")
        }

    }

