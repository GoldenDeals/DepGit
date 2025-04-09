package errors

var (
	//      ERR NAME		  message 		   source	  list of parent errors
	// ERR_NOT_FOUND = New("example error", "errtable", io.EOF, io.EOF, io.EOF)

	ErrNotFound      = New("not found")
	ErrBadData       = New("bad input data")
	ErrAlreadyExists = New("already exists")
)
