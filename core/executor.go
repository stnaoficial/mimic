package core

type Executor struct {
	scanner   *Scanner
	filesRead EntryMap

	generator      *Generator
	filesGenerated EntryMap

	writer       *Writer
	writtenFiles EntryMap
}

func NewExecutor(config *Config) *Executor {
	return &Executor{
		scanner:   NewScanner(config),
		filesRead: make(EntryMap),

		generator:      NewGenerator(config),
		filesGenerated: make(EntryMap),

		writer:       NewWriter(config),
		writtenFiles: make(EntryMap),
	}
}

func (e *Executor) Scan() {
	e.filesRead = e.scanner.Scan()
}

func (e *Executor) Generate() {
	e.filesGenerated = e.generator.Generate(e.filesRead)
}

func (e *Executor) Write() {
	e.writtenFiles = e.writer.Write(e.filesGenerated)
}
