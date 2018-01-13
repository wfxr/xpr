/*******************************************************************************
 *   Author: Wenxuan
 *    Email: wenxuangm@gmail.com
 *  Created: 2018-01-13 10:49
 *******************************************************************************/

package xprogress

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Event func(p *Progress)

func OutputSimpleMessage(p *Progress) {
	fmt.Printf("prev=%d, curr=%d, total=%d, pct=%.2f%%, speed=%.2f/s, eta=%v, eda=%s\n",
		p.Previous(),
		p.Current(),
		p.Total(),
		p.Percent(),
		p.Speed(),
		p.ETA(),
		p.EDA().Format("2006-01-02 15:04:05"),
	)
}

type Progress struct {
	curr, prev, total int64
	interval          time.Duration
	startTime         time.Time
	event             Event
}

func New(total int64) *Progress {
	p := Progress{}
	p.SetTotal(total)
	p.SetInterval(time.Second)
	return &p
}

func (p *Progress) SetEvent(event Event) {
	p.event = event
}

func (p *Progress) SetInterval(d time.Duration) {
	p.interval = d
}

func (p *Progress) Finished() bool {
	return p.Current() >= p.Total()
}

func (p *Progress) Current() int64 {
	return atomic.LoadInt64(&p.curr)
}

func (p *Progress) Previous() int64 {
	return atomic.LoadInt64(&p.prev)
}

func (p *Progress) Total() int64 {
	return atomic.LoadInt64(&p.total)
}

func (p *Progress) SetTotal(total int64) {
	atomic.StoreInt64(&p.total, total)
}

func (p *Progress) syncPrevious() {
	atomic.StoreInt64(&p.prev, p.Current())
}

func (p *Progress) Inc() {
	p.Add(1)
}

func (p *Progress) Add(count int64) {
	atomic.AddInt64(&(p.curr), count)
}

func (p *Progress) Percent() float64 {
	return float64(p.Current()) / float64(p.Total()) * 100
}

func (p *Progress) Speed() float64 {
	elapsed := time.Since(p.startTime)
	return float64(p.Current()) / elapsed.Seconds()
}

func (p *Progress) Remain() int64 {
	return p.Total() - p.Current()
}

func (p *Progress) ETA() time.Duration {
	return time.Duration(float64(p.Remain())/p.Speed()) * time.Second
}

func (p *Progress) EDA() time.Time {
	return time.Now().Add(time.Duration(float64(p.Remain())/p.Speed()) * time.Second)
}

func (p *Progress) Start() {
	p.startTime = time.Now()
	ticker := time.NewTicker(p.interval)
	go func() {
		for _ = range ticker.C {
			if p.event != nil {
				p.event(p)
				p.syncPrevious()
			}
			if p.Finished() {
				ticker.Stop()
			}
		}
	}()
}
