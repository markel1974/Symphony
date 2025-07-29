/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package terminal

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/terminal/vt100"
)

// EquipmentFactory is responsible for creating terminal equipment instances by leveraging provided input-output interfaces.
type EquipmentFactory struct {
}

// NewEquipmentFactory creates and returns a new instance of EquipmentFactory.
func NewEquipmentFactory() *EquipmentFactory {
	return &EquipmentFactory{}
}

// Create initializes a new VT100 terminal instance using the provided input/output and enter key.
func (f *EquipmentFactory) Create(_ string, enterKey rune) interfaces.ITerminal {
	return vt100.NewVt100(enterKey)
}
