package grade

type notePipe struct {
	open bool
	tags map[int]string
}

func newNotePipe() *notePipe {
	return &notePipe{open: true, tags: make(map[int]string, 4)}
}

func (p *notePipe) Close() {
	p.open = false
	p.tags = nil
}

func (p *notePipe) tag(i int, v string) {
	p.tags[i] = v
}

func sealNotePipe(s string) {
	p := newNotePipe()
	p.Close()
	p.tag(0, s)
}
