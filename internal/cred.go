package internal

type Store interface {
	Save(c Credentials) error
	Get() (Credentials, error)
}

type Credentials map[string][]byte