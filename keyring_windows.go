package keyring

import (
	"syscall"

	"github.com/danieljoos/wincred"
)

const errNotFound = "Element not found."

type windowsKeychain struct{}

// Get gets a secret from the keyring given a service name and a user.
func (k windowsKeychain) Get(service, username string) (string, error) {
	cred, err := wincred.GetGenericCredential(k.credName(service, username))
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			// error signaled
			i := int(errno)
			if i == 1168 {
				return "", ErrNotFound
			}
		}
		if err.Error() == errNotFound {
			return "", ErrNotFound
		}
		return "", err
	}

	return string(cred.CredentialBlob), nil
}

// Set stores stores user and pass in the keyring under the defined service
// name.
func (k windowsKeychain) Set(service, username, password string) error {
	cred := wincred.NewGenericCredential(k.credName(service, username))
	cred.UserName = username
	cred.CredentialBlob = []byte(password)
	return cred.Write()
}

// Delete deletes a secret, identified by service & user, from the keyring.
func (k windowsKeychain) Delete(service, username string) error {
	cred, err := wincred.GetGenericCredential(k.credName(service, username))
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			// error signaled
			i := int(errno)
			if i == 1168 {
				return ErrNotFound
			}
		}
		if err.Error() == errNotFound {
			return ErrNotFound
		}
		return err
	}

	return cred.Delete()
}

// credName combines service and username to a single string.
func (k windowsKeychain) credName(service, username string) string {
	return service + ":" + username
}

func init() {
	provider = windowsKeychain{}
}
