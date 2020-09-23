package di

import (
	"os"

	"github.com/golobby/container"
	"github.com/oraksil/orakki/internal/domain/engine"
	"github.com/oraksil/orakki/internal/domain/services"
	"github.com/oraksil/orakki/internal/domain/usecases"
	"github.com/oraksil/orakki/internal/presenter/impl"
	"github.com/oraksil/orakki/internal/presenter/mq/handlers"
	"github.com/oraksil/orakki/pkg/utils"
	"github.com/sangwonl/mqrpc"
)

func newServiceConfig() *services.ServiceConfig {
	hostname, _ := os.Hostname()

	return &services.ServiceConfig{
		MqRpcUri:        utils.GetStrEnv("MQRPC_URI", "amqp://oraksil:oraksil@localhost:5672/"),
		MqRpcNamespace:  utils.GetStrEnv("MQRPC_NAMESPACE", "oraksil"),
		MqRpcIdentifier: utils.GetStrEnv("MQRPC_IDENTIFIER", hostname),

		GipanImageFramesIpcPath: utils.GetStrEnv("IPC_IMAGE_FRAMES", "/var/oraksil/ipc/images.ipc"),
		GipanSoundFramesIpcPath: utils.GetStrEnv("IPC_SOUND_FRAMES", "/var/oraksil/ipc/sounds.ipc"),
		GipanKeyInputsIpcPath:   utils.GetStrEnv("IPC_KEY_INPUTS", "/var/oraksil/ipc/keys.ipc"),
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

func newWebRTCSession() services.WebRTCSession {
	return impl.NewWebRTCSession()
}

func newEngineFactory() engine.EngineFactory {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	return impl.NewGameEngineFactory(serviceConf)
}

func newSetupUseCase() *usecases.SetupUseCase {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	var msgService services.MessageService
	container.Make(&msgService)

	var webRTCSession services.WebRTCSession
	container.Make(&webRTCSession)

	var engineFactory engine.EngineFactory
	container.Make(&engineFactory)

	return &usecases.SetupUseCase{
		ServiceConfig:  serviceConf,
		MessageService: msgService,
		WebRTCSession:  webRTCSession,
		EngineFactory:  engineFactory,
	}
}

func newGamingUseCase() *usecases.GamingUseCase {
	var engineFactory engine.EngineFactory
	container.Make(&engineFactory)

	return &usecases.GamingUseCase{
		EngineFactory: engineFactory,
	}
}

func newSetupHandler() *handlers.SetupHandler {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	var setupUseCase *usecases.SetupUseCase
	container.Make(&setupUseCase)

	return &handlers.SetupHandler{
		ServiceConfig: serviceConf,
		SetupUseCase:  setupUseCase,
	}
}

func newGamingHandler() *handlers.GamingHandler {
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	var gamingUseCase *usecases.GamingUseCase
	container.Make(&gamingUseCase)

	return &handlers.GamingHandler{
		ServiceConfig: serviceConf,
		GamingUseCase: gamingUseCase,
	}
}
