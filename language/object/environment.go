package object

type Module (map[string]Object)

// Program environment
type ProgramEnvironment struct {
	modules      map[string]Module
	RunDirectory string
}

func NewProgramEnvironment(runDirectory string) *ProgramEnvironment {
	return &ProgramEnvironment{
		modules:      make(map[string]Module),
		RunDirectory: runDirectory,
	}
}

func (environment *ProgramEnvironment) IsModuleEvaluated(filepath string) bool {
	return environment.modules[filepath] != nil
}

func (environment *ProgramEnvironment) RegisterModule(filepath string) {
	environment.modules[filepath] = make(Module)
}

func (environment *ProgramEnvironment) RegisterModuleExport(filepath string, name string, value Object) {
	environment.modules[filepath][name] = value
}

// Environment
type Environment struct {
	filepath           string
	store              map[string]Object
	outer              *Environment
	ProgramEnvironment *ProgramEnvironment
}

func NewEnvironment(filepath string, programEnvironment *ProgramEnvironment) *Environment {
	return &Environment{
		filepath:           filepath,
		store:              make(map[string]Object),
		outer:              nil,
		ProgramEnvironment: programEnvironment,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	environment := NewEnvironment(outer.filepath, outer.ProgramEnvironment)
	environment.outer = outer
	return environment
}

func (environment *Environment) Get(name string) (Object, bool) {
	object, found := environment.store[name]

	// Reach for outer variables
	if !found && environment.outer != nil {
		object, found = environment.outer.Get(name)
	}

	return object, found
}

func (environment *Environment) Set(name string, value Object) Object {
	environment.store[name] = value
	return value
}

func (environment *Environment) Export(name string, value Object) {
	environment.ProgramEnvironment.RegisterModuleExport(environment.filepath, name, value)
}

func (environment *Environment) GetModuleValue(filepath string, name string) (Object, bool) {
	value, ok := environment.ProgramEnvironment.modules[filepath][name]
	return value, ok
}
