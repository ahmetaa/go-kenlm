// Package kenlm provides Go bindings to KenLM for scoring text under an
// n-gram language model loaded from a .binary or .arpa file.
package kenlm

/*
#cgo CXXFLAGS: -std=c++17 -O3 -DNDEBUG -DKENLM_MAX_ORDER=6
#cgo CPPFLAGS: -I${SRCDIR} -I${SRCDIR}/third_party/kenlm

#cgo darwin CPPFLAGS: -I/opt/homebrew/include
#cgo darwin LDFLAGS: -L${SRCDIR}/third_party/kenlm/build/lib -lkenlm -lkenlm_util -L/opt/homebrew/lib -lz -lbz2 -llzma

#cgo linux LDFLAGS: -L${SRCDIR}/third_party/kenlm/build/lib -lkenlm -lkenlm_util -lz -lbz2 -llzma -lpthread -lstdc++

#include <stdlib.h>
#include "wrap.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

// Model is a loaded KenLM language model. Safe for concurrent reads.
type Model struct {
	handle *C.kenlm_model
}

// LoadModel loads a KenLM model from an ARPA or binary file at path.
func LoadModel(path string) (*Model, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var cerr *C.char
	h := C.kenlm_load(cpath, &cerr)
	if h == nil {
		msg := "kenlm: failed to load model"
		if cerr != nil {
			msg = "kenlm: " + C.GoString(cerr)
			C.kenlm_string_free(cerr)
		}
		return nil, errors.New(msg)
	}
	m := &Model{handle: h}
	runtime.SetFinalizer(m, func(x *Model) { x.Close() })
	return m, nil
}

// Close releases the underlying model. Safe to call multiple times.
func (m *Model) Close() {
	if m == nil || m.handle == nil {
		return
	}
	C.kenlm_free(m.handle)
	m.handle = nil
	runtime.SetFinalizer(m, nil)
}

// Order returns the n-gram order of the model (e.g. 6 for a 6-gram LM).
func (m *Model) Order() int {
	return int(C.kenlm_order(m.handle))
}

// Score returns the total log10 probability of the whitespace-separated
// token sequence in tokens, including the implicit <s>...</s> markers.
func (m *Model) Score(tokens string) float64 {
	ct := C.CString(tokens)
	defer C.free(unsafe.Pointer(ct))
	return float64(C.kenlm_score_sentence(m.handle, ct))
}
