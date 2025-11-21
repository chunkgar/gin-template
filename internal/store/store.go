package store

var client Factory

// Factory defines the iam platform storage interface.
type Factory interface {
	Close() error
	UserAccount() UserAccount
	User() User
	AdminUser() AdminUser
}

// Client return the store client instance.
func Client() Factory {
	return client
}

// SetClient set the store client.
func SetClient(factory Factory) {
	client = factory
}
