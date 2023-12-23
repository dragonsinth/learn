package main

type blocksByZ [512]map[block]bool

func newBlocksByZ(blocks []block) blocksByZ {
	ret := blocksByZ{}
	for i := range ret {
		ret[i] = map[block]bool{}
	}
	for _, b := range blocks {
		ret.add(b)
	}
	return ret
}

func (bbz blocksByZ) add(b block) {
	for z := b.min(Z); z <= b.max(Z); z++ {
		bbz[z][b] = true
	}
}

func (bbz blocksByZ) rem(b block) {
	for z := b.min(Z); z <= b.max(Z); z++ {
		delete(bbz[z], b)
	}
}
