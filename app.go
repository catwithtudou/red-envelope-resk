package red_envelope

import (
	_ "github.com/catwithtudou/red-envelope-account/core/accounts"
	"github.com/catwithtudou/red-envelope-infra"
	"github.com/catwithtudou/red-envelope-infra/base"
	"github.com/catwithtudou/red-envelope-resk/apis/gorpc"
	_ "github.com/catwithtudou/red-envelope-resk/apis/gorpc"
	_ "github.com/catwithtudou/red-envelope-resk/apis/web"
	_ "github.com/catwithtudou/red-envelope-resk/core/envelopes"
	"github.com/catwithtudou/red-envelope-resk/jobs"
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
