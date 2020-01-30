// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { kdTree } from './kd'

describe('kdTree search', function () {

    const points = [
        { x: 5, y: 5 },
        { x: 2, y: 2 },
        { x: 6, y: 6 },
    ];

    const distance = function (a, b) {
        return Math.pow(a.x - b.x, 2) + Math.pow(a.y - b.y, 2);
    }

    const tree = new kdTree(points, distance, ['x', 'y']);

    it('finds the closest point', function () {
        // Nearby
        const nearest = tree.nearest({ x: 2.1, y: 2.1 }, 1);
        assert.deepEqual(nearest, { x: 2, y: 2 });
    });

    it('finds the closest point even if the distance is zero', function () {
        // Exact match.
        const nearest = tree.nearest({ x: 5, y: 5 }, 1);
        assert.deepEqual(nearest, { x: 5, y: 5 });
    });

    it('finds the closest point even if it is on the other side of the hyperplane', function () {
        // This tests checking down the alternate tree.
        // I.e. the point we are searching for has an x value of 4.5, so it
        // should search down the 2,2 side of the tree, but it will turn out
        // that the hyperplane through 5,5 will be closer than 2,2 so we need to
        // look for a closer match in the 6,6 area, which will actually be our
        // closest point.
        const nearest = tree.nearest({ x: 4.5, y: 7 }, 1);
        assert.deepEqual(nearest, { x: 6, y: 6 });
    });
});