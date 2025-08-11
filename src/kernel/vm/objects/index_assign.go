package objects

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/errors"
)

func IndexAssign(dst IObject, src IObject, selectors []IObject) error {
	numSel := len(selectors)
	for sIdx := numSel - 1; sIdx > 0; sIdx-- {
		next, err := dst.IndexGet(selectors[sIdx])
		if err != nil {
			if errors.Is(err, errors.ErrNotIndexable) {
				return fmt.Errorf("not indexable: %s", dst.TypeName())
			}
			if errors.Is(err, errors.ErrInvalidIndexType) {
				return fmt.Errorf("invalid index type: %s",
					selectors[sIdx].TypeName())
			}
			return err
		}
		dst = next
	}
	if err := dst.IndexSet(selectors[0], src); err != nil {
		if errors.Is(err, errors.ErrNotIndexAssignable) {
			return fmt.Errorf("not index-assignable: %s", dst.TypeName())
		}
		if errors.Is(err, errors.ErrInvalidIndexValueType) {
			return fmt.Errorf("invaid index values type: %s", src.TypeName())
		}
		return err
	}
	return nil
}
