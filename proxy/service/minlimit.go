package service

type MinLimit struct {
	count   int
	curTime int64
}

func (p *MinLimit) Count(now int64) (curCount int) {
	if now-p.curTime > 60 {
		p.count = 1
		p.curTime = now
		curCount = p.count
	}
	p.count++
	curCount = p.count
	return
}

func (p *MinLimit) Check(now int64) int {
	if now-p.curTime > 60 {
		return 0
	}
	return p.count
}
