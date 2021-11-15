package nodes

import "github.com/micro/micro/v3/service/errors"

// ErrFileNotFound returns a 404 error
func ErrFileNotFound(format string, a ...interface{}) error {
	return errors.NotFound("file.not.found", format, a...)
}

// ErrBranchInfoMissing returns a 500 error
func ErrBranchInfoMissing(identifier string) error {
	return errors.InternalServerError("branch.info.incomplete", "Cannot find client for branch %s - did you forget to insert a middleware?", identifier)
}

// ErrBranchInfoRootMissing returns a 500 error
func ErrBranchInfoRootMissing(identifier string) error {
	return errors.InternalServerError("branch.info.incomplete", "Cannot find Root in branch %s - did you forget to insert a middleware?", identifier)
}

// ErrCannotReadStore returns a 403 error
func ErrCannotReadStore(store string) error {
	return errors.Forbidden("forbidden.store", "You are not allowed to access store %s", store)
}

// ErrCannotWriteStore returns a 403 error
func ErrCannotWriteStore(store string) error {
	return errors.Forbidden("forbidden.store", "You are not allowed to write to store %s", store)
}

// ErrCannotFindACL returns a 500 error
func ErrCannotFindACL() error {
	return errors.InternalServerError("branch.info.incomplete", "Cannot load ACLs in branch - did you forget to insert a middleware?")
}
