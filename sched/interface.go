package sched

type jobOps interface {
	DefineNewJob()
	CreateNewJob()
}
