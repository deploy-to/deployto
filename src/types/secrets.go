package types

type Secrets []string

func NewSecrets() Secrets {
	newSecrets := make([]string, 0)
	return newSecrets
}

func (ss *Secrets) Add(s string) {
	*ss = append(*ss, s)
}

//TODO mask logs/dump
