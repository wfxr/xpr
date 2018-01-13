/*******************************************************************************
 *   Author: Wenxuan
 *    Email: wenxuangm@gmail.com
 *  Created: 2018-01-13 11:39
 *******************************************************************************/

package xprogress

import (
	"sync"
	"testing"
)

func equalsInt64(expected, actual int64, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected {%d} was {%d}", expected, actual)
	}
}

func Test_Increment(t *testing.T) {
	p := New(1000 * 1000)
	p.Inc()
	equalsInt64(1, p.Current(), t)
}

func Test_Add(t *testing.T) {
	p := New(9999)
	p.Add(100)
	equalsInt64(100, p.Current(), t)
}

func Test_Increment_ThreadSafe(t *testing.T) {
	p := New(1000 * 1000)
	count := 10 * 1000
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Inc()
		}()
	}
	wg.Wait()
	equalsInt64(int64(count), p.Current(), t)
}

func Test_Add_ThreadSafe(t *testing.T) {
	p := New(1000 * 1000)
	count := 10 * 1000
	wg := sync.WaitGroup{}
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Add(2)
		}()
	}
	wg.Wait()
	equalsInt64(int64(count*2), p.Current(), t)
}
