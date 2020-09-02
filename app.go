package red_envelope

import (
	_ "github.com/catwithtudou/red-envelope-account/core/accounts"
	"github.com/catwithtudou/red-envelope-infra"
	"github.com/catwithtudou/red-envelope-infra/base"
	"red-envelope/apis/gorpc"
	_ "red-envelope/apis/gorpc"
	_ "red-envelope/apis/web"
	_ "red-envelope/core/envelopes"
	"red-envelope/jobs"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.GoRpcStarter{})
	infra.Register(&gorpc.GoRpcApiStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.EurekaStarter{})
	infra.Register(&base.HookStarter{})
}
