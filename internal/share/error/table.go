package errors

var (
	//      ERR NAME		  message 		   source	  list of parent errors
	// ERR_NOT_FOUND = New("example error", "errtable", io.EOF, io.EOF, io.EOF)

	ERR_NOT_FOUND      = New("not found")
	ERR_BAD_DATA       = New("bad input data")
	ERR_ALREADY_EXISTS = New("already exists")
)
