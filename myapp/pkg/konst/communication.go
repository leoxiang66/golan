package konst

var BoolChan chan bool


func init()  {
	BoolChan = make(chan bool,1)
}
