package charger

type Context struct {
	records map[string]*keyRecord
}

func (ctx *Context) render() error {
	deps := make(map[string][]string, len(ctx.records))

	for name, record := range ctx.records {
		depList, err := record.Key.GetRenderer().Dependencies(record.Value)
		if err != nil {
			return err
		}

		deps[name] = depList
	}

	fulfilled := make(map[string]struct{}, len(ctx.records))

	for {
		if len(fulfilled) == len(ctx.records) {
			return nil
		}

		var removed bool

	DepLoop:
		for name, depList := range deps {
			if _, ok := fulfilled[name]; ok {
				continue
			}

			for _, dep := range depList {
				if _, ok := fulfilled[dep]; !ok {
					continue DepLoop
				}
			}

			fulfilled[name] = struct{}{}
			removed = true
		}

		if !removed {
			return ErrCyrcleDetected
		}
	}
}
