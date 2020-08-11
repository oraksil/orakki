package di

import (
	"os"

	"github.com/golobby/container"
	"github.com/sangwonl/mqrpc"
	"gitlab.com/oraksil/orakki/internal/domain/services"
	"gitlab.com/oraksil/orakki/internal/domain/usecases"
	"gitlab.com/oraksil/orakki/internal/presenter/mq/handlers"
	"gitlab.com/oraksil/orakki/pkg/utils"
)

func newServiceConfig() *services.ServiceConfig {
	useStaticOrakki := utils.GetBoolEnv("USE_STATIC_ORAKKI", false)
	var orakkiId string
	if useStaticOrakki {
		orakkiId = utils.GetStrEnv("STATIC_ORAKKI_ID", "orakki-static")
	} else {
		orakkiId, _ = os.Hostname()
	}

	return &services.ServiceConfig{
		MqRpcUri:       utils.GetStrEnv("MQRPC_URI", "amqp://oraksil:oraksil@localhost:5672/"),
		MqRpcNamespace: utils.GetStrEnv("MQRPC_NAMESPACE", "oraksil"),

		UseStaticOrakki: useStaticOrakki,
		OrakkiId:        orakkiId,
		PeerName:        utils.GetStrEnv("PEER_NAME", orakkiId),
	}
}

func newMqService() *mqrpc.MqService {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	svc, err := mqrpc.NewMqService(serviceConf.MqRpcUri, serviceConf.MqRpcNamespace)
	if err != nil {
		panic(err)
	}
	return svc
}

func newMessageService() services.MessageService {
	var mqService *mqrpc.MqService
	container.Make(&mqService)

	return &mqrpc.DefaultMessageServiceImpl{MqService: mqService}
}

func newSystemMonitorUseCase() *usecases.SystemStateMonitorUseCase {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	var msgService services.MessageService
	container.Make(&msgService)

	return &usecases.SystemStateMonitorUseCase{
		ServiceConfig:  serviceConf,
		MessageService: msgService,
	}
}

func newSystemHandler() *handlers.SystemHandler {
	var sysMonUseCase *usecases.SystemStateMonitorUseCase
	container.Make(&sysMonUseCase)

	return &handlers.SystemHandler{
		SystemMonitorUseCase: sysMonUseCase,
	}
}
