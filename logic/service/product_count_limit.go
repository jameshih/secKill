package service

import "sync"

// small case for private varible
type ProductCountMgr struct {
	productCount map[int]int
	lock         sync.RWMutex
}

func NewProductCountMgr() (productCountMgr *ProductCountMgr) {
	productCountMgr = &ProductCountMgr{
		productCount: make(map[int]int, 128),
	}
	return
}

func (p *ProductCountMgr) Count(productID int) (count int) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	count, _ = p.productCount[productID]
	return
}

func (p *ProductCountMgr) Add(productID, count int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	cur, ok := p.productCount[productID]
	if !ok {
		cur = count
	} else {
		cur += count
	}
	p.productCount[productID] = cur
}
