package red_envelope

import (
	"red-envelope/infra"
	"red-envelope/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	//infra.Register(&base.DbxDatabaseStarter{})
	//infra.Register(&base.ValidatorStarter{})
	//infra.Register(&base.IrisServerStarter{})

}
