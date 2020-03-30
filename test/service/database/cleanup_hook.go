package database

type cleanupHook struct {
	hookName         string
	attachedEntities []entity
}

func newCleanUpHook(name string) *cleanupHook {
	return &cleanupHook{hookName: name}
}

func (ch *cleanupHook) pushEntity(e entity) {
	ch.attachedEntities = append(ch.attachedEntities, e)
}

func (ch *cleanupHook) name() string {
	return ch.hookName
}

func (ch *cleanupHook) entities() []entity {
	return ch.attachedEntities
}
