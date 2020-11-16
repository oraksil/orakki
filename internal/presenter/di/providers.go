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

		GipanImageFramesIpcUri: utils.GetStrEnv("IPC_IMAGE_FRAMES", "tcp://127.0.0.1:8765"),
		GipanSoundFramesIpcUri: utils.GetStrEnv("IPC_SOUND_FRAMES", "tcp://127.0.0.1:8766"),
		GipanCmdInputsIpcUri:   utils.GetStrEnv("IPC_CMD_INPUTS", "tcp://127.0.0.1:8767"),

		TurnServerUri:      utils.GetStrEnv("TURN_URI", ""),
		TurnServerUsername: utils.GetStrEnv("TURN_USERNAME", ""),
		TurnServerPassword: utils.GetStrEnv("TURN_PASSWORD", ""),
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
	var serviceConf *services.ServiceConfig
	container.Make(&serviceConf)

	return impl.NewWebRTCSession(
		serviceConf.TurnServerUri,
		serviceConf.TurnServerUsername,
		serviceConf.TurnServerPassword,
	)
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
	var msgService services.MessageService
	container.Make(&msgService)

	var engineFactory engine.EngineFactory
	container.Make(&engineFactory)

	return &usecases.GamingUseCase{
		MessageService: msgService,
		EngineFactory:  engineFactory,
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
