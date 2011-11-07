package procfs

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
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

	Uptime    int
	Idletime  int
}

const (
	PROCFS_PROCESSES = "Processes"
	PROCFS_SELF = "Self"
	PROCFS_UPTIME = "Uptime"
	PROCFS_IDLETIME = "Idletime"
)

func (pfs *ProcFS) Fill() {
	pfs.List(PROCFS_PROCESSES)
	for _, p := range pfs.Processes {
		p.Fill()
	}
	pfs.Get(PROCFS_SELF)

	pfs.Get(PROCFS_UPTIME)
	pfs.Get(PROCFS_IDLETIME)
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
	var uf = path.Join(procfsdir, "uptime")
	switch k {
	case PROCFS_SELF:
		var selfdir = path.Join(procfsdir, "self")
		if !exists(selfdir) {
			return
		}
		fi, _ := os.Readlink(selfdir)
		pfs.Self = fi
	case PROCFS_UPTIME:
		str, err := ioutil.ReadFile(uf)
		if err == nil {
			ss := strings.Split(string(str), " ")
			if len(ss) >= 2 {
				pfs.Uptime, _ = strconv.Atoi(ss[0])
			}
		}
	case PROCFS_IDLETIME:
		str, err := ioutil.ReadFile(uf)
		if err == nil {
			ss := strings.Split(string(str), " ")
			if len(ss) >= 2 {
				pfs.Idletime, _ = strconv.Atoi(ss[1])
			}
		}
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
	switch k {
	case PROCFS_PROC_AUXV:
		p.Auxv, _ = ioutil.ReadFile(path.Join(pdir, "auxv"))
	case PROCFS_PROC_CMDLINE:
		cl, err := ioutil.ReadFile(path.Join(pdir, "cmdline"))
		if err == nil {
			p.Cmdline = splitNull(cl)
		}
	case PROCFS_PROC_CWD:
		p.Cwd, _ = os.Readlink(path.Join(pdir, "cwd"))
	case PROCFS_PROC_ENVIRON:
		envB, err := ioutil.ReadFile(path.Join(pdir, "environ"))
		if err == nil {
			p.Environ = make(map[string]string)
			envS := splitNull(envB)
			for _, s := range envS {
				// split on =
				ss := strings.SplitN(s, "=", 2)
				if len(ss) == 2 {
					p.Environ[ss[0]] = ss[1]
				}
			}
		}
	case PROCFS_PROC_EXE:
		p.Exe, _ = os.Readlink(path.Join(pdir, "exe"))
	case PROCFS_PROC_ROOT:
		p.Root, _ = os.Readlink(path.Join(pdir, "root"))
	case PROCFS_PROC_STATUS:
		statLines, err := ioutil.ReadFile(path.Join(pdir, "status"))
		if err == nil {
			p.Status = make(map[string]string)
			statS := strings.Split(string(statLines), "\n")
			for _, s := range statS {
				ss := strings.SplitN(s, ":", 2)
				if len(ss) == 2 {
					p.Status[ss[0]] = strings.TrimSpace(ss[1])
				}
			}
		}
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
