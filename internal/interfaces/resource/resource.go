package resource

type Entity any
type Response map[string]any

type ResourceInterface[P Entity, R Response] interface {
	Resource(entity *P) R
	Collection(entities []*P) []R
}
