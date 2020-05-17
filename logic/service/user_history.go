package service

import "sync"

type UserBuyHistory struct {
	History map[int]int
	Lock    sync.RWMutex
}

func (p *UserBuyHistory) GetProductBuyCount(productID int) (count int) {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	count, _ = p.History[productID]
	return
}

func (p *UserBuyHistory) Add(productID, count int) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	cur, ok := p.History[productID]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	p.History[productID] = cur
}
