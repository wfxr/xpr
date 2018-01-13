/*******************************************************************************
 *   Author: Wenxuan
 *    Email: wenxuangm@gmail.com
 *  Created: 2018-01-13 12:44
 *******************************************************************************/

package main

import (
	"sync"
	"time"

	"github.com/wfxr/xprogress"
)

func main() {
	p := xprogress.New(1000 * 1000)
	p.SetEvent(xprogress.OutputSimpleMessage)
	p.SetInterval(time.Millisecond * 500)
	p.Start()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				p.Inc()
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}
	wg.Wait()
}
