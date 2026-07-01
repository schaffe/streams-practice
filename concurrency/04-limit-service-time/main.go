//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var userQuota sync.Map = sync.Map{}

const quotaConst int32 = 2

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {

	ctx := context.Background()
	if !u.IsPremium {
		quota := getQuota(u.ID)
		if quota.Load() <= 0 {
			return false
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		ctx = timeoutCtx

		//ticker to update counter
		tick := time.Tick(time.Second)
		go func() {
			for {
				select {
				case <-tick:
					q := getQuota(u.ID)
					println(q.Load())
					if q.Add(-1) <= 0 {
						fmt.Printf("User %d user up their quota, aborting\n", u.ID)
						cancel()
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()

	}

	ch := make(chan struct{})

	go func() {
		defer close(ch)
		process()
	}()

	select {
	case <-ch:
		return true
	case <-ctx.Done():
		return false
	}
}

func getQuota(userID int) *atomic.Int32 {
	if q, ok := userQuota.Load(userID); ok {
		return q.(*atomic.Int32)
	}

	var quota atomic.Int32
	quota.Store(quotaConst + 1)

	quotaAny, _ := userQuota.LoadOrStore(userID, &quota)
	quotaPtr := quotaAny.(*atomic.Int32)
	return quotaPtr
}

func main() {
	RunMockServer()
}
