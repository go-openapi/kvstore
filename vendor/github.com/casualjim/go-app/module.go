package app

// LifecycleCallback function definition
type LifecycleCallback interface {
	Call(Application) error
}

// Init is an initializer for an initalization function
type Init func(Application) error

// Call implements the callback interface
func (fn Init) Call(app Application) error {
	return fn(app)
}

// Start is an initializer for a start function
type Start func(Application) error

// Call implements the callback interface
func (fn Start) Call(app Application) error {
	return fn(app)
}

// Stop is an initializer for a stop function
type Stop func(Application) error

// Call implements the callback interface
func (fn Stop) Call(app Application) error {
	return fn(app)
}

// Reload is an initalizater for a reload function
type Reload func(Application) error

// Call implements the callback interface
func (fn Reload) Call(app Application) error {
	return fn(app)
}

// A Module is a component that has a specific lifecycle
type Module interface {
	Init(Application) error
	Start(Application) error
	Stop(Application) error
	Reload(Application) error
}

// MakeModule by passing the callback functions.
// You can pass multiple callback functions of the same type if you want
func MakeModule(callbacks ...LifecycleCallback) Module {
	var (
		init   []Init
		start  []Start
		reload []Reload
		stop   []Stop
	)

	for _, callback := range callbacks {
		switch cb := callback.(type) {
		case Init:
			init = append(init, cb)
		case Start:
			start = append(start, cb)
		case Stop:
			stop = append(stop, cb)
		case Reload:
			reload = append(reload, cb)
		}
	}

	return &dynamicModule{
		init:   init,
		start:  start,
		reload: reload,
		stop:   stop,
	}
}

type dynamicModule struct {
	init   []Init
	start  []Start
	stop   []Stop
	reload []Reload
}

func (d *dynamicModule) Init(app Application) error {
	for _, cb := range d.init {
		if err := cb.Call(app); err != nil {
			return err
		}
	}
	return nil
}

func (d *dynamicModule) Start(app Application) error {
	for _, cb := range d.start {
		if err := cb.Call(app); err != nil {
			return err
		}
	}
	return nil
}

func (d *dynamicModule) Stop(app Application) error {
	for _, cb := range d.stop {
		if err := cb.Call(app); err != nil {
			return err
		}
	}
	return nil
}

func (d *dynamicModule) Reload(app Application) error {
	for _, cb := range d.reload {
		if err := cb.Call(app); err != nil {
			return err
		}
	}
	return nil
}
