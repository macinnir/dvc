package gates

import "github.com/macinnir/dvc/core/lib/utils/types"

const ROOT_USER_ID int64 = 1

func AccountGate(model types.IAccountable, actor types.IUserContainer) bool {
	return model.Account() == actor.Account() || actor.ID() == ROOT_USER_ID
}
