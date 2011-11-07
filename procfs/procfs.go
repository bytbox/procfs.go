package procfs

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const procfsdir = "/proc"

type Filler interface {
	Fill()
}

type Lister interface {
	List(string)
}

type Getter interface {
	Get(string)
}

type ProcFS struct {
	Processes map[string]*Process
	Self      string
}

const (
	PROCFS_PROCESSES = "Processes"
	PROCFS_SELF = "Self"
)

func (pfs *ProcFS) Fill() {
	pfs.List(PROCFS_PROCESSES)
	pfs.Get(PROCFS_SELF)
}

func (pfs *ProcFS) List(k string) {
	switch k {
	case PROCFS_PROCESSES:
		if !exists(procfsdir) {
			return
		}
		pfs.Processes = make(map[string]*Process)
		ds, err := ioutil.ReadDir(procfsdir)
		if err != nil {
			return
		}
		// get all numeric entries
		for _, d := range ds {
			n := d.Name
			id, err := strconv.Atoi(n)
			if isNumeric(n) && err == nil {
				proc := Process{PID: id}
				pfs.Processes[n] = &proc
			}
		}
	}
}

func (pfs *ProcFS) Get(k string) {
	switch k {
	case PROCFS_SELF:
		var selfdir = path.Join(procfsdir, "self")
		if !exists(selfdir) {
			return
		}
		fi, _ := os.Readlink(selfdir)
		pfs.Self = fi
	}
}

type Process struct {
	PID     int
	Auxv    []byte
	Cmdline []string
	Cwd     string
	Environ map[string]string
	Exe     string
	Fds     map[string]*Fd
	Root    string
	Status  map[string]string
	Threads map[string]*Thread
}
// TODO limits, maps, mem, mountinfo, mounts, mountstats, ns, smaps, stat

const (
	PROCFS_PROC_AUXV = "Process.Auxv"
	PROCFS_PROC_CMDLINE = "Process.Cmdline"
	PROCFS_PROC_CWD = "Process.Cwd"
	PROCFS_PROC_ENVIRON = "Process.Environ"
	PROCFS_PROC_EXE = "Process.Exe"
	PROCFS_PROC_ROOT = "Process.Root"
	PROCFS_PROC_STATUS = "Process.Status"

	PROCFS_PROC_FDS = "Process.Fds"
	PROCFS_PROC_THREADS = "Process.Threads"
)

func (p *Process) Fill() {
	p.Get(PROCFS_PROC_AUXV)
	p.Get(PROCFS_PROC_CMDLINE)
	p.Get(PROCFS_PROC_CWD)
	p.Get(PROCFS_PROC_ENVIRON)
	p.Get(PROCFS_PROC_EXE)
	p.Get(PROCFS_PROC_ROOT)
	p.Get(PROCFS_PROC_STATUS)

	// Fds
	p.List(PROCFS_PROC_FDS)
	for _, f := range p.Fds {
		f.Fill()
	}

	// Threads
	p.List(PROCFS_PROC_THREADS)
	for _, t := range p.Threads {
		t.Fill()
	}
}

func (p *Process) List(k string) {

}

func (p *Process) Get(k string) {
	pdir := path.Join(procfsdir, strconv.Itoa(p.PID))
	println(pdir)
	switch k {
	case PROCFS_PROC_AUXV:
		p.Auxv, _ = ioutil.ReadFile(path.Join(pdir, "auxv"))
	case PROCFS_PROC_CMDLINE:
		cl, err := ioutil.ReadFile(path.Join(pdir, "cmdline"))
		if err == nil {
			p.Cmdline = splitNull(cl)
		}
		println(err.Error())
	case PROCFS_PROC_CWD:
		p.Cwd, _ = os.Readlink(path.Join(pdir, "cwd"))
	case PROCFS_PROC_ENVIRON:
	}
}

type Fd struct {
	Path  string
	Pos   int
	Flags int
}

const (
	PROCFS_PROC_FD_PATH = "Process.Fd.Path"
	PROCFS_PROC_FD_POS = "Process.Fd.Pos"
	PROCFS_PROC_FD_FLAGS = "Process.Fd.Flags"
)

func (f *Fd) Fill() {
	f.Get(PROCFS_PROC_FD_PATH)
	f.Get(PROCFS_PROC_FD_POS)
	f.Get(PROCFS_PROC_FD_FLAGS)
}

func (f *Fd) Get(k string) {
	switch k {

	}
}

type Thread struct {
	// TODO
}

func (t *Thread) Fill() {

}

func (t *Thread) Get(k string) {

}
