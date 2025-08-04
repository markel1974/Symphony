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

package interfaces

// ISurface represents a drawing interface for rendering and manipulating graphical or textual content.
// GetSize retrieves the dimensions of the surface as rows and columns.
// Draw places a character at the specified position on the surface.
// DrawColor places a character with specified foreground, background colors, and color mode.
// DrawText renders a string starting at the specified position on the surface.
// DrawTextColor renders a string with specified foreground, background colors, and color mode.
// DrawSeries draws a series of data points within the given dimensions and value range.
type ISurface interface {
	GetSize() (int, int)

	Draw(rows int, column int, c rune)

	DrawColor(rows int, column int, c rune, fg ColorDef, bg ColorDef, mode ColorMode)

	DrawText(rows int, column int, c string)

	DrawTextColor(rows int, column int, c string, fg ColorDef, bg ColorDef, mode ColorMode)

	DrawSeries(data []float64, w int, h int, min float64, max float64)

	Begin()

	End()
}
