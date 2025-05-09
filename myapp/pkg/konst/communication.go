package konst

var BoolChan chan bool
var StringChan chan string


func init()  {
	BoolChan = make(chan bool,1)
	StringChan = make(chan string,1)
}
