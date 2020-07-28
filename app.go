package red_envelope

import (
	_ "red-envelope/apis/web"
	_ "red-envelope/core/accounts"
	_ "red-envelope/core/envelopes"
	"red-envelope/infra"
	"red-envelope/infra/base"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})

}
