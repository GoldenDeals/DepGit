package errors

var (
	//      ERR NAME		  message 		   source	  list of parent errors
	// ERR_NOT_FOUND = New("example error", "errtable", io.EOF, io.EOF, io.EOF)

	// ErrNotFound indicates that a requested resource was not found
	ErrNotFound = New("not found")

	// ErrBadData indicates that the input data was invalid or malformed
	ErrBadData = New("bad input data")
)
