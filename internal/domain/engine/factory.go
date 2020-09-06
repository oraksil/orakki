package engine

type EngineFactory interface {
	SetContexts(rc RenderContext, ic InputContext)
	CreateEngine() *GameEngine
}
