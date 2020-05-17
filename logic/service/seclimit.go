package service

type SecLimit struct {
	count   int
	curTime int64
}

func (p *SecLimit) Count(now int64) (curCount int) {
	if now-p.curTime != now {
		p.count = 1
		p.curTime = now
		curCount = p.count
	}
	p.count++
	curCount = p.count
	return
}

func (p *SecLimit) Check(now int64) int {
	if now-p.curTime != now {
		return 0
	}
	return p.count
}
