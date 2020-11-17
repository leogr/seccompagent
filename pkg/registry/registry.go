package registry

import (
	specs "github.com/opencontainers/runtime-spec/specs-go"
	libseccomp "github.com/seccomp/libseccomp-golang"
	"golang.org/x/sys/unix"
)

type HandlerResult struct {
	Intr   bool
	ErrVal int32
	Val    uint64
	Flags  uint32
}

type HandlerFunc func(libseccomp.ScmpFd, *libseccomp.ScmpNotifReq) HandlerResult

// Helper functions for handlers
func HandlerResultIntr() HandlerResult {
	return HandlerResult{Intr: true}
}
func HandlerResultContinue() HandlerResult {
	return HandlerResult{Flags: libseccomp.NotifRespFlagContinue}
}

func HandlerResultErrno(err error) HandlerResult {
	if err == nil {
		return HandlerResult{}
	}
	errno, ok := err.(unix.Errno)
	if !ok {
		return HandlerResult{ErrVal: int32(unix.ENOSYS), Val: ^uint64(0)}
	}
	if errno == 0 {
		return HandlerResult{}
	}
	return HandlerResult{ErrVal: int32(errno), Val: ^uint64(0)}
}
func HandlerResultSuccess() HandlerResult {
	return HandlerResult{}
}

// Registry

type Registry struct {
	SyscallHandler map[string]HandlerFunc
}

type ResolverFunc func(state *specs.ContainerProcessState) *Registry

func New() *Registry {
	return &Registry{
		SyscallHandler: map[string]HandlerFunc{},
	}
}

func (r *Registry) Add(name string, f HandlerFunc) {
	r.SyscallHandler[name] = f
}