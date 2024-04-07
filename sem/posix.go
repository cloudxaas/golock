// +build linux darwin

package posixsem

/*
#include <fcntl.h>
#include <sys/stat.h>
#include <stdlib.h> 
#include <semaphore.h>
#include <errno.h>

sem_t *sem_open_wrapper(const char *name, int oflag, mode_t mode, unsigned int value) {
    return sem_open(name, oflag, mode, value);
}
*/
import "C"
import (
    "errors"
    "unsafe"
)

// Sem represents a named semaphore.
type Sem struct {
    name *C.char
    sem  *C.sem_t
}

// Open opens a named semaphore.
func Open(name string, value uint) (*Sem, error) {
    cName := C.CString(name)
    defer C.free(unsafe.Pointer(cName))

    // Remove O_EXCL flag to allow opening an existing semaphore.
    sem := C.sem_open_wrapper(cName, C.O_CREAT, C.S_IRUSR|C.S_IWUSR, C.uint(value))
    if sem == C.SEM_FAILED {
        return nil, errors.New("failed to open semaphore")
    }
    return &Sem{name: cName, sem: sem}, nil
}

// Wait decreases the semaphore value (lock/wait).
func (s *Sem) Wait() error {
    if C.sem_wait(s.sem) == -1 {
        return errors.New("failed to wait on semaphore")
    }
    return nil
}



// Post increases the semaphore value (unlock/post).
func (s *Sem) Post() error {
    if C.sem_post(s.sem) == -1 {
        return errors.New("failed to post semaphore")
    }
    return nil
}

// Close closes the semaphore.
func (s *Sem) Close() error {
    if C.sem_close(s.sem) == -1 {
        return errors.New("failed to close semaphore")
    }
    return nil
}

// Unlink removes a named semaphore.
func Unlink(name string) error {
    cName := C.CString(name)
    defer C.free(unsafe.Pointer(cName))
    
    // Attempt to unlink the semaphore.
    if C.sem_unlink(cName) == -1 {
        return errors.New("failed to unlink semaphore")
    }
    return nil
}
