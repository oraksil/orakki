package engine

type EngineFactory interface {
	SetContexts(rc RenderContext, ic InputContext)
	CanCreateEngine() bool
	CreateEngine() *GameEngine
}
