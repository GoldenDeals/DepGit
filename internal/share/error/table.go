package error

var (
	//      ERR NAME		  message 		   source	  list of parent errors
	// ERR_NOT_FOUND = New("example error", "errtable", io.EOF, io.EOF, io.EOF)

	ERR_NOT_FOUND = New("not found")
)
