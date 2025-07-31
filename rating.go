// Copyright 2025 Robert Snedegar
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigo

import "fmt"

// Rating pairs up a BigO and the score it received in the current processing.
type Rating struct {
	bigO  *BigO
	score float64
}

func (r *Rating) String() string {
	// TODO(rsned): Figure out how many significant digits are useful.
	return fmt.Sprintf("%s: %0.3f", r.bigO.label, r.score)
}

// BigO returns the BigO instance associated with this rating.
func (r *Rating) BigO() *BigO {
	return r.bigO
}

// Score returns the correlation score for this rating.
func (r *Rating) Score() float64 {
	return r.score
}

// defaultRating is used when nothing has been processed yet.
var defaultRating = &Rating{
	bigO:  defaultBigO,
	score: 0,
}
