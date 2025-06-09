package transprocessor

/*
	transprocessor implements the transactions structure that whole project complies with
*/

/* permissions processor */
type PermProcessor struct {
	errCh		chan<-error
}
