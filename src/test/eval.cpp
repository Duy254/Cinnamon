/*
    Cinnamon UCI chess engine
    Copyright (C) Giuseppe Cannella

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

#if defined(DEBUG_MODE) || defined(FULL_TEST)

#include <gtest/gtest.h>
#include <set>
#include "../SearchManager.h"

TEST(eval, eval1) {
    SearchManager &searchManager = Singleton<SearchManager>::getInstance();
    searchManager.loadFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
    int score = searchManager.getScore(WHITE, false);
    EXPECT_EQ(5, score);
    score = searchManager.getScore(BLACK, false);
    EXPECT_EQ(-5, score);
}

//TEST(eval, passed_pawn) {TODO
//    SearchManager &searchManager = Singleton<SearchManager>::getInstance();
//    searchManager.loadFen("1q2k3/4r3/2P1n2p/2b2Qp1/2N5/8/1B3PPP/3R2K1 w - - 0 40");
//    int score = searchManager.getScore(WHITE, false);
//    EXPECT_GT(score, 400);
//    EXPECT_LE(score, 600);
//}
#endif