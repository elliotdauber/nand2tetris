package main

type SymbolTable struct {
	table      map[string]SymbolTableEntry
	var_counts map[SymbolKind]int
}

func (this *SymbolTable) Init() SymbolTable {
	this.table = make(map[string]SymbolTableEntry)
	this.var_counts = map[SymbolKind]int{
		SymbolStatic: 0,
		SymbolField:  0,
		SymbolArg:    0,
		SymbolVar:    0,
	}
	return *this
}

func (this *SymbolTable) Reset() {
	for k := range this.table {
		delete(this.table, k)
	}

	for k := range this.var_counts {
		this.var_counts[k] = 0
	}
}

func (this *SymbolTable) Define(name string, typ string, kind SymbolKind) {
	index := this.VarCount(kind)
	this.increment_var_count(kind)
	this.table[name] = SymbolTableEntry{
		typ:   typ,
		kind:  kind,
		index: index,
	}
}

func (this *SymbolTable) VarCount(kind SymbolKind) int {
	if val, found := this.var_counts[kind]; found {
		return val
	}
	return -1
}

func (this *SymbolTable) KindOf(name string) SymbolKind {
	if val, found := this.table[name]; found {
		return val.kind
	}
	return SymbolNone
}

func (this *SymbolTable) TypeOf(name string) string {
	if val, found := this.table[name]; found {
		return val.typ
	}
	return ""
}

func (this *SymbolTable) IndexOf(name string) int {
	if val, found := this.table[name]; found {
		return val.index
	}
	return -1
}

func (this *SymbolTable) FindSymbol(name string) (SymbolTableEntry, bool) {
	val, found := this.table[name]
	return val, found
}

func (this *SymbolTable) increment_var_count(kind SymbolKind) {
	if val, found := this.var_counts[kind]; found {
		this.var_counts[kind] = val + 1
	}
}

type SymbolTableEntry struct {
	typ   string
	kind  SymbolKind
	index int
}

type SymbolKind int

const (
	SymbolStatic SymbolKind = iota
	SymbolField
	SymbolArg
	SymbolVar
	SymbolNone
)
