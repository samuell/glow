package glow

var baseConv = [256]byte{
	'A':  'T',
	'T':  'A',
	'C':  'G',
	'G':  'C',
	'N':  'N',
	'\n': '\n',
}

type BaseComplementer struct {
	In  chan []byte
	Out chan []byte
}

func (self *BaseComplementer) Init() {
	go func() {
		for line := range self.In {
			if line[0] != '>' {
				for pos := range line {
					line[pos] = baseConv[line[pos]]
				}
			}
			self.Out <- append([]byte(nil), line...)
		}
		close(self.Out)
	}()
}
