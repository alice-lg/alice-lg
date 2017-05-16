package sources

type Source interface {
	Status()
	Neighbours()
	Routes(neighbourId int)
}
