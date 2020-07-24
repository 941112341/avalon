package repository

type IdGeneratorRepository interface {
	Get() (*IdGenerator, error)
	Save(generator IdGenerator) error
	UpdateVersion(generator IdGenerator) (int64, error)
	FindByMaxIDBetween(left, right int64) (*IdGenerator, error)
}
