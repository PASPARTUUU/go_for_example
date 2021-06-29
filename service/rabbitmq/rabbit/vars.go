package rabbit

const (
	// --- Exchanges

	ExchangeUser = "e_user"

	// --- Queues

	// user
	QueueNewUser     = "q_new_user"
	QueueUpdatedUser = "q_updated_user"
	QueueDeletedUser = "q_deleted_user"

	// --- Routing keys

	// common
	KeyNew    = "k_new"
	KeyUpdate = "k_update"
	KeyDelete = "k_delete"

	//
	KeyUser = "k_user"

	// --- Consumers

	ConsumerNewUser     = "c_new_user"
	ConsumerUpdatedUser = "c_updated_user"
	ConsumerDeletedUser = "c_deleted_user"
)
