package di

import (
	"os"

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

		GipanImageFramesIpcUri: utils.GetStrEnv("IPC_IMAGE_FRAMES", "ipc://../gipan/images.ipc"),
		GipanSoundFramesIpcUri: utils.GetStrEnv("IPC_SOUND_FRAMES", "tcp://../gipan/sounds.ipc"),
		GipanCmdInputsIpcUri:   utils.GetStrEnv("IPC_CMD_INPUTS", "tcp://../gipan/cmds/ipc"),

		TurnServerUri:       utils.GetStrEnv("TURN_URI", ""),
		TurnServerSecretKey: utils.GetStrEnv("TURN_SECRET_KEY", ""),
		TurnServerTTL:       utils.GetIntEnv("TURN_TTL", 3600),

		PlayerHealthCheckTimeout:  utils.GetIntEnv("PLAYER_HEALTHCHECK_TIMEOUT", 20),
		PlayerHealthCheckInterval: utils.GetIntEnv("PLAYER_HEALTHCHECK_INTERVAL", 3),
		PlayerIdleCheckTimeout:    utils.GetIntEnv("PLAYER_IDLECHECK_TIMEOUT", 60),
		PlayerIdleCheckInterval:   utils.GetIntEnv("PLAYER_IDLECHECK_INTERVAL", 7),
	}
}

func newMqService(serviceConf *services.ServiceConfig) *mqrpc.MqService {
	svc, err := mqrpc.NewMqService(serviceConf.MqRpcUri, serviceConf.MqRpcNamespace, serviceConf.MqRpcIdentifier)
	if err != nil {
		panic(err)
	}
	return svc
}

func newMessageService(mqService *mqrpc.MqService) services.MessageService {
	return &mqrpc.DefaultMessageService{MqService: mqService}
}

func newWebRTCSession(serviceConf *services.ServiceConfig) services.WebRTCSession {
	return impl.NewWebRTCSession(
		serviceConf.TurnServerUri,
		serviceConf.TurnServerSecretKey,
		serviceConf.TurnServerTTL,
	)
}

func newGipanDriver(serviceConf *services.ServiceConfig) engine.GipanDriver {
	return impl.NewGipanDriver(
		serviceConf.GipanImageFramesIpcUri,
		serviceConf.GipanSoundFramesIpcUri,
		serviceConf.GipanCmdInputsIpcUri,
	)
}

func newEngineFactory(
	serviceConf *services.ServiceConfig,
	gipanDrv engine.GipanDriver) engine.EngineFactory {
	return impl.NewGameEngineFactory(serviceConf, gipanDrv)
}

func newSetupUseCase(
	serviceConf *services.ServiceConfig,
	msgService services.MessageService,
	webRTCSession services.WebRTCSession,
	engineFactory engine.EngineFactory) *usecases.SetupUseCase {
	return &usecases.SetupUseCase{
		ServiceConfig:  serviceConf,
		MessageService: msgService,
		WebRTCSession:  webRTCSession,
		EngineFactory:  engineFactory,
	}
}

func newGamingUseCase(
	serviceConf *services.ServiceConfig,
	msgService services.MessageService,
	engineFactory engine.EngineFactory,
	gipanDrv engine.GipanDriver) *usecases.GamingUseCase {
	return &usecases.GamingUseCase{
		ServiceConfig:  serviceConf,
		MessageService: msgService,
		EngineFactory:  engineFactory,
		GipanDriver:    gipanDrv,
	}
}

func newSetupHandler(
	serviceConf *services.ServiceConfig,
	setupUseCase *usecases.SetupUseCase) *handlers.SetupHandler {
	return &handlers.SetupHandler{
		ServiceConfig: serviceConf,
		SetupUseCase:  setupUseCase,
	}
}

func newGamingHandler(
	serviceConf *services.ServiceConfig,
	gamingUseCase *usecases.GamingUseCase) *handlers.GamingHandler {
	return &handlers.GamingHandler{
		ServiceConfig: serviceConf,
		GamingUseCase: gamingUseCase,
	}
}
