package actor

type Actor struct {
	ID      string
	ClerkID string
}

func NewActor(id string, clerkID string) *Actor {
	return &Actor{ID: id, ClerkID: clerkID}
}

func (a *Actor) GetID() string {
	return a.ID
}
