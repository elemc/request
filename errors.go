package request

type ErrInitializedAlready bool
const errInitializedAlready = ErrInitializedAlready(true)
func (err ErrInitializedAlready) Error() string {
	return "the request package has been initialized already"
}
