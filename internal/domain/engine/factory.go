package engine

type EngineFactory interface {
	SetContexts(rc RenderContext, ic InputContext, sc SessionContext)
	CanCreateEngine() bool
	CreateEngine() *GameEngine
}
