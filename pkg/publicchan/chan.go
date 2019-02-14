package publicchan

var ChanTask chan string

func init() {
	ChanTask = make(chan string, 4096)
}
