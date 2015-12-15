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

#pragma once

#include "GenMoves.h"
#include "Endgame.h"
#include <fstream>
#include <string.h>
#include <iomanip>

using namespace _board;

class Eval : public GenMoves {

public:


    Eval();

    virtual ~Eval();

    int getScore(const int side, const int alpha = -_INFINITE, const int beta = _INFINITE, const bool print = false);

    template<int side>
    int lazyEval() {
        return lazyEvalSide<side>() - lazyEvalSide<side ^ 1>();
    }

#ifdef DEBUG_MODE
    unsigned lazyEvalCuts;
#endif
protected:
    STATIC_CONST int FUTIL_MARGIN = 154;        //CLOP 144
    STATIC_CONST int EXT_FUTILY_MARGIN = 392;    //CLOP 405
    STATIC_CONST int RAZOR_MARGIN = 1071;        //CLOP 1070
    STATIC_CONST int ATTACK_KING = 30;
    STATIC_CONST int BISHOP_ON_QUEEN = 2;
    STATIC_CONST int BACKWARD_PAWN = 2;
    STATIC_CONST int NO_PAWNS = 15;
    STATIC_CONST int DOUBLED_ISOLATED_PAWNS = 14;
    STATIC_CONST int DOUBLED_PAWNS = 5;
    STATIC_CONST int ENEMIES_PAWNS_ALL = 8;
    STATIC_CONST int PAWN_7H = 32;
    STATIC_CONST int PAWN_CENTER = 15;
    STATIC_CONST int PAWN_IN_RACE = 114;
    STATIC_CONST int PAWN_ISOLATED = 3;
    STATIC_CONST int PAWN_NEAR_KING = 2;
    STATIC_CONST int PAWN_BLOCKED = 5;
    STATIC_CONST int UNPROTECTED_PAWNS = 5;
    STATIC_CONST int ENEMY_NEAR_KING = 2;
    STATIC_CONST int FRIEND_NEAR_KING = 1;
    STATIC_CONST int BISHOP_NEAR_KING = 10;
    STATIC_CONST int HALF_OPEN_FILE = 3;
    STATIC_CONST int KNIGHT_TRAPPED = 5;
    STATIC_CONST int END_OPENING = 6;
    STATIC_CONST int BONUS2BISHOP = 18;
    STATIC_CONST int CONNECTED_ROOKS = 7;
    STATIC_CONST int ROOK_OPEN_FILE = 10;
    STATIC_CONST int OPEN_DIAG = 10;
    STATIC_CONST int ROOK_SEMI_OPEN_FILE = 5;
    STATIC_CONST int ROOK_7TH_RANK = 10;
    STATIC_CONST int ROOK_BLOCKED = 13;
    STATIC_CONST int ROOK_TRAPPED = 6;
    STATIC_CONST int UNDEVELOPED = 4;
    STATIC_CONST int UNDEVELOPED_BISHOP = 4;
#ifdef DEBUG_MODE
    typedef struct {
        int BAD_BISHOP[2];
        int MOB_BISHOP[2];
        int UNDEVELOPED_BISHOP[2];
        int OPEN_DIAG_BISHOP[2];
        int BONUS2BISHOP[2];

        int MOB_PAWNS[2];
        int ATTACK_KING_PAWN[2];
        int PAWN_CENTER[2];
        int PAWN_7H[2];
        int PAWN_IN_RACE[2];
        int PAWN_BLOCKED[2];
        int UNPROTECTED_PAWNS[2];
        int PAWN_ISOLATED[2];
        int DOUBLED_PAWNS[2];
        int DOUBLED_ISOLATED_PAWNS[2];
        int BACKWARD_PAWN[2];
        int FORK_SCORE[2];
        int PAWN_PASSED[2];
        int ENEMIES_PAWNS_ALL[2];
        int NO_PAWNS[2];

        int KING_SECURITY_BISHOP[2];
        int KING_SECURITY_QUEEN[2];
        int KING_SECURITY_KNIGHT[2];
        int KING_SECURITY_ROOK[2];
        int DISTANCE_KING[2];
        int END_OPENING_KING[2];
        int PAWN_NEAR_KING[2];
        int MOB_KING[2];

        int MOB_QUEEN[2];

        int BISHOP_ON_QUEEN[2];
        int HALF_OPEN_FILE[2];

        int UNDEVELOPED_KNIGHT[2];
        int KNIGHT_TRAPPED[2];
        int MOB_KNIGHT[2];


        int ROOK_7TH_RANK[2];
        int ROOK_TRAPPED[2];
        int MOB_ROOK[2];
        int ROOK_BLOCKED[2];
        int ROOK_OPEN_FILE[2];
        int ROOK_SEMI_OPEN_FILE[2];
        int CONNECTED_ROOKS[2];
    } _TSCORE_DEBUG;
    _TSCORE_DEBUG SCORE_DEBUG;
#endif

private:
    enum _Tstatus {
        OPEN, MIDDLE, END
    };

    const int MOB_QUEEN[3][29] = {{0,   1,   1,   1, 1, 1, 1, 1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1,  1},
                                  {-10, -9,  -5,  0, 3, 6, 7, 10, 11, 12, 15, 18, 28, 30, 32, 35, 40, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61},
                                  {-20, -15, -10, 0, 1, 3, 4, 9,  11, 12, 15, 18, 28, 30, 32, 33, 34, 36, 37, 39, 40, 41, 42, 43, 44, 45, 56, 47, 48}};

    const int MOB_ROOK[3][15] = {{-1,  0,   1,  4, 5, 6,  7,  9,  12, 14, 19, 22, 23, 24, 25},
                                 {-9,  -8,  1,  8, 9, 10, 15, 20, 28, 30, 40, 45, 50, 51, 52},
                                 {-15, -10, -5, 0, 9, 11, 16, 22, 30, 32, 40, 45, 50, 51, 52}};

    const int MOB_KNIGHT[9] = {-8, -4, 7, 10, 15, 20, 30, 35, 40};

    const int MOB_BISHOP[3][14] = {{-8,  -7,  2,  8, 9, 10, 15, 20, 28, 30, 40, 45, 50, 50},
                                   {-20, -10, -4, 0, 5, 10, 15, 20, 28, 30, 40, 45, 50, 50},
                                   {-20, -10, -4, 0, 3, 8,  13, 18, 25, 30, 40, 45, 50, 50}};

    const int MOB_KING[3][9] = {{1,   2,   2,   1,  0,  0,  0,  0,  0},
                                {-5,  0,   5,   5,  5,  0,  0,  0,  0},
                                {-50, -30, -10, 10, 25, 40, 50, 55, 60}};

    const int MOB_CASTLE[3][3] = {{-50, 30, 50},
                                  {-1,  10, 10},
                                  {0,   0,  0}};

    const int MOB_PAWNS[34] = {-1, 2, 3, 4, 5, 10, 12, 14, 18, 22, 25, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 70, 75, 80, 90, 95, 100, 110};

    const int BONUS_ATTACK_KING[18] = {-1, 2, 8, 64, 128, 512, 512, 512, 512, 512, 512, 512, 512, 512, 512, 512, 512, 512};

    const u64 PAWN_PROTECTED_MASK[2][64] = {{0x200ULL, 0x500ULL, 0xa00ULL, 0x1400ULL, 0x2800ULL, 0x5000ULL, 0xa000ULL, 0x4000ULL, 0x20000ULL, 0x50000ULL, 0xa0000ULL, 0x140000ULL, 0x280000ULL, 0x500000ULL, 0xa00000ULL, 0x400000ULL, 0x2000000ULL, 0x5000000ULL, 0xa000000ULL, 0x14000000ULL, 0x28000000ULL, 0x50000000ULL, 0xa0000000ULL, 0x40000000ULL, 0x200000000ULL, 0x500000000ULL, 0xa00000000ULL, 0x1400000000ULL, 0x2800000000ULL, 0x5000000000ULL, 0xa000000000ULL, 0x4000000000ULL, 0x20000000000ULL, 0x50000000000ULL, 0xa0000000000ULL, 0x140000000000ULL, 0x280000000000ULL, 0x500000000000ULL, 0xa00000000000ULL, 0x400000000000ULL, 0x2000000000000ULL, 0x5000000000000ULL, 0xa000000000000ULL, 0x14000000000000ULL, 0x28000000000000ULL, 0x50000000000000ULL, 0xa0000000000000ULL, 0x40000000000000ULL, 0xFF000000000000ULL, 0xFF000000000000ULL, 0xFF000000000000ULL, 0xFF000000000000ULL, 0xFF000000000000ULL, 0xFF000000000000ULL, 0xFF000000000000ULL, 0xFF000000000000ULL, 0,                     0,                     0,                     0,                     0,                     0,                     0,                     0},
                                            {0,        0,        0,        0,         0,         0,         0,         0,         0xFF00ULL,  0xFF00ULL,  0xFF00ULL,  0xFF00ULL,   0xFF00ULL,   0xFF00ULL,   0xFF00ULL,   0xFF00ULL,   0x200ULL,     0x500ULL,     0xa00ULL,     0x1400ULL,     0x2800ULL,     0x5000ULL,     0xa000ULL,     0x4000ULL,     0x20000ULL,     0x50000ULL,     0xa0000ULL,     0x140000ULL,     0x280000ULL,     0x500000ULL,     0xa00000ULL,     0x400000ULL,     0x2000000ULL,     0x5000000ULL,     0xa000000ULL,     0x14000000ULL,     0x28000000ULL,     0x50000000ULL,     0xa0000000ULL,     0x40000000ULL,     0x200000000ULL,     0x500000000ULL,     0xa00000000ULL,     0x1400000000ULL,     0x2800000000ULL,     0x5000000000ULL,     0xa000000000ULL,     0x4000000000ULL,     0x20000000000ULL,    0x50000000000ULL,    0xa0000000000ULL,    0x140000000000ULL,   0x280000000000ULL,   0x500000000000ULL,   0xa00000000000ULL,   0x400000000000ULL,   0xFF00000000000000ULL, 0xFF00000000000000ULL, 0xFF00000000000000ULL, 0xFF00000000000000ULL, 0xFF00000000000000ULL, 0xFF00000000000000ULL, 0xFF00000000000000ULL, 0xFF00000000000000ULL}};

    const u64 PAWN_BACKWARD_MASK[2][64] = {{0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 131584ULL, 328960ULL, 657920ULL, 1315840ULL, 2631680ULL, 5263360ULL, 10526720ULL, 4210688ULL, 33685504ULL, 84213760ULL, 168427520ULL, 336855040ULL, 673710080ULL, 1347420160ULL, 2694840320ULL, 1077936128ULL, 8623489024ULL, 21558722560ULL, 43117445120ULL, 86234890240ULL, 172469780480ULL, 344939560960ULL, 689879121920ULL, 275951648768ULL, 2207613190144ULL, 5519032975360ULL, 11038065950720ULL, 22076131901440ULL, 44152263802880ULL, 88304527605760ULL, 176609055211520ULL, 70643622084608ULL, 565148976676864ULL, 1412872441692160ULL, 2825744883384320ULL, 5651489766768640ULL, 11302979533537280ULL, 22605959067074560ULL, 45211918134149120ULL, 18084767253659648ULL, 562949953421312ULL, 1407374883553280ULL, 2814749767106560ULL, 5629499534213120ULL, 11258999068426240ULL, 22517998136852480ULL, 45035996273704960ULL, 18014398509481984ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0},
                                           {0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 512ULL,    1280ULL,   2560ULL,   5120ULL,    10240ULL,   20480ULL,   40960ULL,    16384ULL,   131584ULL,   328960ULL,   657920ULL,    1315840ULL,   2631680ULL,   5263360ULL,    10526720ULL,   4210688ULL,    33685504ULL,   84213760ULL,    168427520ULL,   336855040ULL,   673710080ULL,    1347420160ULL,   2694840320ULL,   1077936128ULL,   8623489024ULL,    21558722560ULL,   43117445120ULL,    86234890240ULL,    172469780480ULL,   344939560960ULL,   689879121920ULL,    275951648768ULL,   2207613190144ULL,   5519032975360ULL,    11038065950720ULL,   22076131901440ULL,   44152263802880ULL,    88304527605760ULL,    176609055211520ULL,   70643622084608ULL,    565148976676864ULL, 1412872441692160ULL, 2825744883384320ULL, 5651489766768640ULL, 11302979533537280ULL, 22605959067074560ULL, 45211918134149120ULL, 18084767253659648ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0ULL, 0}};

    const u64 PAWN_PASSED_MASK[2][64] = {{0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x3ULL,               0x7ULL,               0xeULL,               0x1cULL,               0x38ULL,               0x70ULL,               0xe0ULL,               0xc0ULL,               0x303ULL,             0x707ULL,             0xe0eULL,             0x1c1cULL,             0x3838ULL,             0x7070ULL,             0xe0e0ULL,             0xc0c0ULL,             0x30303ULL,           0x70707ULL,           0xe0e0eULL,           0x1c1c1cULL,           0x383838ULL,           0x707070ULL,           0xe0e0e0ULL,           0xc0c0c0ULL,           0x3030303ULL,         0x7070707ULL,         0xe0e0e0eULL,         0x1c1c1c1cULL,         0x38383838ULL,         0x70707070ULL,         0xe0e0e0e0ULL,         0xc0c0c0c0ULL,         0x303030303ULL,       0x707070707ULL,       0xe0e0e0e0eULL,       0x1c1c1c1c1cULL,       0x3838383838ULL,       0x7070707070ULL,       0xe0e0e0e0e0ULL,       0xc0c0c0c0c0ULL,       0x30303030303ULL,     0x70707070707ULL,     0xe0e0e0e0e0eULL,     0x1c1c1c1c1c1cULL,     0x383838383838ULL,     0x707070707070ULL,     0xe0e0e0e0e0e0ULL,     0xc0c0c0c0c0c0ULL,     0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0},
                                         {0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x303030303030000ULL, 0x707070707070000ULL, 0xe0e0e0e0e0e0000ULL, 0x1c1c1c1c1c1c0000ULL, 0x3838383838380000ULL, 0x7070707070700000ULL, 0xe0e0e0e0e0e00000ULL, 0xc0c0c0c0c0c00000ULL, 0x303030303000000ULL, 0x707070707000000ULL, 0xe0e0e0e0e000000ULL, 0x1c1c1c1c1c000000ULL, 0x3838383838000000ULL, 0x7070707070000000ULL, 0xe0e0e0e0e0000000ULL, 0xc0c0c0c0c0000000ULL, 0x303030300000000ULL, 0x707070700000000ULL, 0xe0e0e0e00000000ULL, 0x1c1c1c1c00000000ULL, 0x3838383800000000ULL, 0x7070707000000000ULL, 0xe0e0e0e000000000ULL, 0xc0c0c0c000000000ULL, 0x303030000000000ULL, 0x707070000000000ULL, 0xe0e0e0000000000ULL, 0x1c1c1c0000000000ULL, 0x3838380000000000ULL, 0x7070700000000000ULL, 0xe0e0e00000000000ULL, 0xc0c0c00000000000ULL, 0x303000000000000ULL, 0x707000000000000ULL, 0xe0e000000000000ULL, 0x1c1c000000000000ULL, 0x3838000000000000ULL, 0x7070000000000000ULL, 0xe0e0000000000000ULL, 0xc0c0000000000000ULL, 0x300000000000000ULL, 0x700000000000000ULL, 0xe00000000000000ULL, 0x1c00000000000000ULL, 0x3800000000000000ULL, 0x7000000000000000ULL, 0xe000000000000000ULL, 0xc000000000000000ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0ULL, 0x0}};

    const uchar PAWN_PASSED[2][64] = {{200, 200, 200, 200, 200, 200, 200, 200, 100, 100, 100, 100, 100, 100, 100, 100, 40, 40, 40, 40, 40, 40, 40, 40, 19, 19, 19, 21, 21, 19, 19, 19, 13, 13, 13, 25, 25, 13, 13, 13, 0,  0,  0,  0,  0,  0,  0,  0,  0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0},
                                      {0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,   0,  0,  0,  0,  0,  0,  0,  0,  13, 13, 13, 25, 25, 13, 13, 13, 19, 19, 19, 21, 21, 19, 19, 19, 40, 40, 40, 40, 40, 40, 40, 40, 100, 100, 100, 100, 100, 100, 100, 100, 200, 200, 200, 200, 200, 200, 200, 200}};

    const u64 PAWN_ISOLATED_MASK[64] = {0x202020202020202ULL, 0x505050505050505ULL, 0xA0A0A0A0A0A0A0AULL, 0x1414141414141414ULL, 0x2828282828282828ULL, 0x5050505050505050ULL, 0xA0A0A0A0A0A0A0A0ULL, 0x4040404040404040ULL, 0x202020202020202ULL, 0x505050505050505ULL, 0xA0A0A0A0A0A0A0AULL, 0x1414141414141414ULL, 0x2828282828282828ULL, 0x5050505050505050ULL, 0xA0A0A0A0A0A0A0A0ULL, 0x4040404040404040ULL, 0x202020202020202ULL, 0x505050505050505ULL, 0xA0A0A0A0A0A0A0AULL, 0x1414141414141414ULL, 0x2828282828282828ULL, 0x5050505050505050ULL, 0xA0A0A0A0A0A0A0A0ULL, 0x4040404040404040ULL, 0x202020202020202ULL, 0x505050505050505ULL, 0xA0A0A0A0A0A0A0AULL, 0x1414141414141414ULL, 0x2828282828282828ULL, 0x5050505050505050ULL, 0xA0A0A0A0A0A0A0A0ULL, 0x4040404040404040ULL, 0x202020202020202ULL, 0x505050505050505ULL, 0xA0A0A0A0A0A0A0AULL, 0x1414141414141414ULL, 0x2828282828282828ULL, 0x5050505050505050ULL, 0xA0A0A0A0A0A0A0A0ULL, 0x4040404040404040ULL, 0x202020202020202ULL, 0x505050505050505ULL,
                                        0xA0A0A0A0A0A0A0AULL, 0x1414141414141414ULL, 0x2828282828282828ULL, 0x5050505050505050ULL, 0xA0A0A0A0A0A0A0A0ULL,
                                        0x4040404040404040ULL, 0x202020202020202ULL, 0x505050505050505ULL, 0xA0A0A0A0A0A0A0AULL, 0x1414141414141414ULL, 0x2828282828282828ULL, 0x5050505050505050ULL, 0xA0A0A0A0A0A0A0A0ULL, 0x4040404040404040ULL, 0x202020202020202ULL, 0x505050505050505ULL, 0xA0A0A0A0A0A0A0AULL, 0x1414141414141414ULL, 0x2828282828282828ULL, 0x5050505050505050ULL, 0xA0A0A0A0A0A0A0A0ULL, 0x4040404040404040ULL};

    const char DISTANCE_KING_OPENING[64] = {-8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -12, -12, -12, -12, -8, -8, -8, -8, -12, -16, -16, -12, -8, -8, -8, -8, -12, -16, -16, -12, -8, -8, -8, -8, -12, -12, -12, -12, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8, -8};

    const char DISTANCE_KING_ENDING[64] = {12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 16, 16, 16, 16, 12, 12, 12, 12, 16, 20, 20, 16, 12, 12, 12, 12, 16, 20, 20, 16, 12, 12, 12, 12, 16, 16, 16, 16, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12};

    Bits &bits = Bits::getInstance();
#ifdef DEBUG_MODE
    int evaluationCount[2];
#endif

    void openFile();

    template<int side, _Tstatus status>
    int evaluatePawn();

    template<int side, _Tstatus status>
    int evaluateBishop(const u64, u64);

    template<int side, _Tstatus status>
    int evaluateQueen(u64 enemies, u64 friends);

    template<int side, _Tstatus status>
    int evaluateKnight(const u64, const u64);

    template<int side, Eval::_Tstatus status>
    int evaluateRook(const u64, u64 enemies, u64 friends);

    template<_Tstatus status>
    int evaluateKing(int side, u64 squares);

    template<int side>
    int lazyEvalSide() {
        return Bits::bitCount(chessboard[PAWN_BLACK + side]) * VALUEPAWN + Bits::bitCount(chessboard[ROOK_BLACK + side]) * VALUEROOK + Bits::bitCount(chessboard[BISHOP_BLACK + side]) * VALUEBISHOP + Bits::bitCount(chessboard[KNIGHT_BLACK + side]) * VALUEKNIGHT + Bits::bitCount(chessboard[QUEEN_BLACK + side]) * VALUEQUEEN;
    }
};

