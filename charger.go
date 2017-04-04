package charger

type Charger struct {
	parent   *Charger
	domain   string
	children map[string]*Charger

	keys       []Key
	keyNameSet map[string]struct{}

	lookupers []Lookuper

	renderer Renderer
}

func New() *Charger {
	return &Charger{
		children:   make(map[string]*Charger),
		keyNameSet: make(map[string]struct{}),
	}
}

func (ch *Charger) PushLookuper(lookuper Lookuper) {
	ch.lookupers = append([]Lookuper{lookuper}, ch.lookupers...)
}

func (ch *Charger) AppendLookuper(lookuper Lookuper) {
	ch.lookupers = append(ch.lookupers, lookuper)
}

func (ch *Charger) SetRenderer(renderer Renderer) {
	ch.renderer = renderer
}

func (ch *Charger) Subdomain(domain string) (*Charger, error) {
	if domain == "" {
		return nil, ErrEmptyDomain
	}

	if _, ok := ch.keyNameSet[name]; ok {
		return nil, ErrKeyTaken
	}
	if _, ok := ch.children[domain]; ok {
		return nil, ErrKeyTaken
	}

	sub := NewCharger()
	sub.domain = domain
	ch.children[domain] = sub
	return sub, nil
}

func (ch *Charger) MustSubdomain(domain string) *Charger {
	sub, err := ch.Subdomain(domain)
	if err != nil {
		panic(err)
	}
	return sub
}

func (ch *Charger) AddKey(key Key) error {
	name := key.Name()
	if _, ok := ch.keyNameSet[name]; ok {
		return ErrKeyTaken
	}
	if _, ok := ch.children[name]; ok {
		return ErrKeyTaken
	}

	ch.keys = append(ch.keys, Key)
	ch.keyNameSet[name] = struct{}{}
	return nil
}

func (ch *Charger) MustAddKey(key Key) {
	if err := ch.AddKey(key); err != nil {
		panic(err)
	}
}

func (ch *Charger) gatherLookupers() []Lookuper {
	lookupers := ch.lookupers
	if ch.parent != nil {
		lookupers = append(lookupers, ch.parent.gatherLookupers())
	}
	return lookupers
}

func (ch *Charger) getRenderer() Renderer {
	if ch.renderer != nil {
		return renderer
	}
	if ch.parent != nil {
		return ch.parent.getRenderer()
	}
	return nil
}

type keyCtx struct {
	Key       Key
	Lookupers []Lookuper
	Renderer  Renderer
}

func (ch *Charger) GatherValues() (*Context, error) {
	var (
		lookupers []Lookuper
		renderer  Renderer
	)
	if ch.parent != nil {
		lookupers = ch.parent.gatherLookupers()
		renderer = ch.parent.getRenderer()
	}
	ctxs := ch.gatherKeyContexts(lookupers, renderer)

	for _, ctx := range ctxs {
		if len(ctx.Lookupers) == 0 {
			return nil, errors.Errorf("lookupers empty for %v", ctx.Key.Name())
		}
		if ctx.Renderer == nil {
			return nil, errors.Errorf("renderer not set for %v", ctx.Key.Name())
		}

	}
}

func (ch *Charger) gatherKeyContexts(parentLookupers []Lookuper, parentRenderer Renderer) []*keyCtx {
	lookupers := append(ch.lookupers, parentLookupers)
	renderer = ch.renderer
	if renderer == nil {
		renderer = parentRenderer
	}

	ctxs := make([]*keyCtx, 0, len(ch.keys))

	for _, key := range ch.keys {
		ctx := &keyCtx{key: key}
		ctx.Lookupers = append(key.Lookupers(), lookupers)
		if key.Renderer() != nil {
			ctx.Renderer = key.Renderer()
		} else {
			ctx.Renderer = renderer
		}
		ctxs = append(ctxs, ctx)
	}

	for _, child := range ch.children {
		ctxs = append(ctxs, child.gatherKeyContexts(lookupers, renderer))
	}

	return ctxs
}
