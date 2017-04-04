package charger

type Charger struct {
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

func (ch *Charger) GatherValues() (*Context, error) {
	ctx := &Context{
		records: make(map[string]*keyRecord, len(ch.records)),
	}

KeyLoop:
	for name, record := range ch.records {
		for _, lookuper := range record.Key.GetLookupers() {
			value, err := lookuper.Lookup(name)
			switch err {
			case nil:
				ctx.records[name] = &keyRecord{record.Key, value}
				continue KeyLoop
			case ErrNotFound:
			default:
				return nil, err
			}
		}

		for _, lookuper := range ch.lookupers {
			value, err := lookuper.Lookup(name)
			switch err {
			case nil:
				ctx.records[name] = &keyRecord{record.Key, value}
				continue KeyLoop
			case ErrNotFound:
			default:
				return nil, err
			}
		}

		value, ok := record.Key.GetDefault()
		if ok {
			ctx.records[name] = &keyRecord{record.Key, value}
			continue KeyLoop
		}

		if record.Key.GetRequired() {
			return nil, &ErrRequired{name}
		}
	}

	if err := ctx.RenderValues(); err != nil {
		return nil, err
	}

	return ctx, nil
}
