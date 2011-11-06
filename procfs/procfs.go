package procfs

type ProcFS struct {
	Processes map[int]Process
	Self      Process
}

type Process struct {
	Auxv    []byte
	Cmdline []string
	Cwd     string
	Environ map[string]string
	Exe     string
	Fds     map[int]Fd
	Root    string
	Status  map[string]string
	Threads map[int]Thread
}

// TODO limits, maps, mem, mountinfo, mounts, mountstats, ns, smaps, stat

type Fd struct {
	Path  string
	Pos   int
	Flags int
}

type Thread struct {
	// TODO
}
