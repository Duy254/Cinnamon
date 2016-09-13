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

package main

const (
	MAX_MOVE int = 130
	RANK_1 uint64 = 0xff00
	RANK_3 uint64 = 0xff000000
	RANK_4 uint64 = 0xff00000000
	RANK_6 uint64 = 0xff000000000000
	STANDARD_MOVE_MASK uint8 = 0x3
	ENPASSANT_MOVE_MASK uint8 = 0x1
	PROMOTION_MOVE_MASK uint8 = 0x2
	MAX_REP_COUNT int = 1024
	NO_PROMOTION int = -1

	TABJUMPPAWN uint64 = 0xFF00000000FF00;
	TABCAPTUREPAWN_RIGHT uint64 = 0xFEFEFEFEFEFEFEFE;
	TABCAPTUREPAWN_LEFT uint64 = 0x7F7F7F7F7F7F7F7F;


)

type  GenMoves struct {
	ChessBoard
	perftMode  bool
	listId     int
	gen_list   []_TmoveP
	//repetitionMap *uint64
	currentPly int
	numMoves   uint64
	numMovesq  uint64
	forceCheck bool
}

func ( self *GenMoves )  setPerft(b bool) {
	self.perftMode = b;
}

func NewGenMoves() *GenMoves {
	p := new(GenMoves)
	p.gen_list = make([]_TmoveP, MAX_PLY)

	return p
}

func ( self *GenMoves ) generateMoves(side uint, allpieces uint64) {

	self.tryAllCastle(side, allpieces);
	self.performDiagShift(BISHOP_BLACK + side, side, allpieces);
	self.performRankFileShift(ROOK_BLACK + side, side, allpieces);
	self.performRankFileShift(QUEEN_BLACK + side, side, allpieces);
	self.performDiagShift(QUEEN_BLACK + side, side, allpieces);
	self.performPawnShift(side, allpieces^0);
	self.performKnightShiftCapture(KNIGHT_BLACK + side, allpieces^0, side);
	self.performKingShiftCapture(side, allpieces^0);
}

func ( self *GenMoves )  generateCaptures(side uint, enemies uint64, friends uint64) bool {

	var allpieces = enemies | friends;
	if self.performPawnCapture(enemies, side) {
		return true;
	}
	if self.performKingShiftCapture(side, enemies) {
		return true;
	}
	if self.performKnightShiftCapture(KNIGHT_BLACK + side, enemies, side) {
		return true;
	}
	if self.performDiagCapture(BISHOP_BLACK + side, enemies, side, allpieces) {
		return true;
	}
	if self.performRankFileCapture(ROOK_BLACK + side, enemies, side, allpieces) {
		return true;
	}
	if self.performRankFileCapture(QUEEN_BLACK + side, enemies, side, allpieces) {
		return true;
	}
	if self.performDiagCapture(QUEEN_BLACK + side, enemies, side, allpieces) {
		return true;
	}
	return false;
}

func ( self *GenMoves ) getForceCheck() bool {
	return self.forceCheck;
}

func ( self *GenMoves ) setForceCheck(b bool) {
	self.forceCheck = b;
}

//func ( self *GenMoves ) getMoveFromSan(fenStr *string, mov *_Tmove) int {

/* self.chessboard [ENPASSANT_IDX] =  NO_ENPASSANT;
 memset(mov, 0, sizeof (_Tmove));
 const string MATCH_QUEENSIDE = "O-O-O e1c1 e8c8";
 const string MATCH_QUEENSIDE_WHITE = "O-O-O e1c1";
 const string MATCH_KINGSIDE_WHITE = "O-O e1g1";
 const string MATCH_QUEENSIDE_BLACK = "O-O-O e8c8";
 const string MATCH_KINGSIDE_BLACK = "O-O e8g8";

 if ((MATCH_QUEENSIDE_WHITE.find(fenStr) != string::npos | |
 MATCH_KINGSIDE_WHITE.find(fenStr) != string::npos) & &
 getPieceAt < WHITE > (POW2[E1]) == KING_WHITE) | |
 ((MATCH_QUEENSIDE_BLACK.find(fenStr) != string::npos | |
 MATCH_KINGSIDE_BLACK.find(fenStr) != string::npos) & &
 getPieceAt < BLACK > (POW2[E8]) == KING_BLACK) {
 if MATCH_QUEENSIDE.find(fenStr) != string::npos {
 mov.typee = QUEEN_SIDE_CASTLE_MOVE_MASK;
 mov.from = QUEEN_SIDE_CASTLE_MOVE_MASK;
 } else {
 mov.from = KING_SIDE_CASTLE_MOVE_MASK;
 mov.typee = KING_SIDE_CASTLE_MOVE_MASK;
 }
 if (fenStr.find("1") != string::npos) {
 mov.side = WHITE;
 } else if (fenStr.find("8") != string::npos) {
 mov.side = BLACK;
 } else {
 panic!();
 }
 mov.from = - 1;
 mov.capturedPiece = SQUARE_FREE;
 return mov.side;
 }
 int from = - 1;
 int to = -1;
 for i 0..64{
 if (!fenStr.compare(0, 2, BOARD[i])) {
 from = i;
 break;
 }
 }
 if from == - 1 {
 cout < < fenStr < < endl;
 panic!();
 }
 for i 0..64{
 if (!fenStr.compare(2, 2, BOARD[i])) {
 to = i;
 break;
 }
 }
 if (to == - 1) {
 cout < < fenStr < < endl;
 panic!();
 }
 int pieceFrom;
 if ((pieceFrom = getPieceAt <WHITE > (POW2[from])) != 12) {
 move.side = WHITE;
 } else if ((pieceFrom = getPieceAt < BLACK >(POW2[from])) != 12) {
 move.side = BLACK;
 } else {
 cout < < "fenStr: " < < fenStr < < " from: " < < from < < endl;
 panic!();
 }
 move.from = from;
 move.to = to;
 if (fenStr.length() == 4) {
 move.typee = STANDARD_MOVE_MASK;
 if (pieceFrom == PAWN_WHITE | | pieceFrom == PAWN_BLACK) {
 if (FILE_AT[from] != FILE_AT[to] & & ( move.side ^ 1 ? getPieceAt < WHITE > (POW2[to]) : getPieceAt < BLACK > (POW2[to])) == SQUARE_FREE) {
 move.typee = ENPASSANT_MOVE_MASK;
 }
 }
 } else if (fenStr.length() == 5) {
 move.typee = PROMOTION_MOVE_MASK;
 if ( move.side == WHITE) {
 move.promotionPiece = INV_FEN[toupper(fenStr.at(4))];
 } else {
 move.promotionPiece = INV_FEN[(uchar) fenStr.at(4)];
 }
 debug_assert!( move.promotionPiece != - 1);
 }
 if ( move.side == WHITE) {
 move.capturedPiece = getPieceAt < BLACK > (POW2[ move.to]);
 move.pieceFrom = getPieceAt < WHITE > (POW2[ move.from]);
 } else {
 move.capturedPiece = getPieceAt < WHITE > (POW2[ move.to]);
 move.pieceFrom = getPieceAt < BLACK > (POW2[ move.from]);
 }
 if ( move.typee == ENPASSANT_MOVE_MASK) {
 move.capturedPiece = ! move.side;
 }
 return move.side;*/
//}

func ( self *GenMoves ) init() {
	self.numMoves = 0;
	self.numMovesq = 0;
	self.listId = 0;
	//        #ifdef DEBUG_MODE
	//        nCutFp = nCutRazor = 0;
	//        betaEfficiency = 0.0;
	//        nCutAB = 0;
	//        nNullMoveCut = 0;
	//        #endif
}

func ( self *GenMoves ) loadFen(fen *string) int {
	//self.repetitionMapCount = 0;
	var side = self.loadFen(fen);
	if side == 2 {
		panic("Bad FEN position format ");
	}
	return side;
}

func ( self *GenMoves ) getDiagCapture(position uint, allpieces uint64, enemies uint64) uint64 {
	return self.getDiagonalAntiDiagonal(position, allpieces) & enemies;
}

func ( self *GenMoves ) getDiagShiftAndCapture(position uint, enemies uint64, allpieces uint64) uint64 {
	var nuovo = self.getDiagonalAntiDiagonal(position, allpieces);
	return (nuovo & enemies) | (nuovo & !allpieces);
}

func ( self *GenMoves ) takeback(mov *_Tmove, oldkey uint64, rep bool) {
	//if rep {
	//	self.popStackMove();
	//}
	self.chessboard[ZOBRISTKEY_IDX] = oldkey;
	self.chessboard[ENPASSANT_IDX] = NO_ENPASSANT;

	var pieceFrom uint;
	var posTo uint;
	var posFrom uint;
	var movecapture uint;
	self.chessboard[RIGHT_CASTLE_IDX] = (mov.typee & 0xf0);
	if (mov.typee & 0x3) == STANDARD_MOVE_MASK || (mov.typee & 0x3) == ENPASSANT_MOVE_MASK {
		posTo = mov.to;
		posFrom = mov.from;
		movecapture = mov.capturedPiece;

		pieceFrom = mov.pieceFrom;
		self.chessboard[pieceFrom] = (self.chessboard[pieceFrom] & NOTPOW2[posTo]) | POW2[posFrom];
		if movecapture != SQUARE_FREE {
			if (mov.typee & 0x3) != ENPASSANT_MOVE_MASK {
				self.chessboard[movecapture] |= POW2[posTo ];
			} else {

				if mov.side != 0 {
					self.chessboard[movecapture] |= POW2[(posTo - 8)];
				} else {
					self.chessboard[movecapture] |= POW2[(posTo + 8)];
				}
			}
		}
	} else if (mov.typee & 0x3) == PROMOTION_MOVE_MASK {
		posTo = mov.to;
		posFrom = mov.from;
		movecapture = mov.capturedPiece;

		self.chessboard[mov.side] |= POW2[posFrom];
		self.chessboard[mov.promotionPiece] &= NOTPOW2[posTo];
		if movecapture != SQUARE_FREE {
			self.chessboard[movecapture] |= POW2[posTo];
		}
	} else if mov.typee & 0xc != 0 {
		//castle
		self.unPerformCastle(mov.side, mov.typee);
	}
}

//func ( self *GenMoves ) setRepetitionMapCount(i uint) {
//	self.repetitionMapCount = i;
//}

func ( self *GenMoves )  getDiagShiftCount(position uint, allpieces uint64) uint {

	return bitCount((self.getDiagonalAntiDiagonal(position, allpieces) & !allpieces));
}

func ( self *GenMoves )  performKingShiftCapture(side uint, enemies uint64) bool {

	var pos = BITScanForward(self.chessboard[KING_BLACK + side]);

	var x1 = enemies & NEAR_MASK1[pos];
	for ; x1 != 0; {
		if self.pushmove(STANDARD_MOVE_MASK, pos, BITScanForward(x1), side, NO_PROMOTION, (KING_BLACK + side)) {
			return true;
		}
		RESET_LSB(x1);
	};
	return false;
}

func ( self *GenMoves )  performKnightShiftCapture(piece uint, enemies uint64, side uint) bool {

	var x = self.chessboard[piece];
	for ; x != 0; {
		var pos = BITScanForward(x);
		var x1 = enemies & KNIGHT_MASK[pos];
		for ; x1 != 0; {
			if self.pushmove(STANDARD_MOVE_MASK, pos, BITScanForward(x1), side, NO_PROMOTION, piece) {
				return true;
			}
			RESET_LSB(x1);
		};
		RESET_LSB(x);
	}
	return false;
}

func ( self *GenMoves ) performDiagCapture(piece uint, enemies uint64, side uint, allpieces uint64) bool {

	var x2 = self.chessboard[piece];
	for ; x2 != 0; {
		var position = BITScanForward(x2);
		var diag = self.getDiagonalAntiDiagonal(position, allpieces) & enemies;

		for ; diag != 0; {
			if self.pushmove(STANDARD_MOVE_MASK, position, BITScanForward(diag), side, NO_PROMOTION, piece) {
				return true;
			}
			RESET_LSB(diag);
		}
		RESET_LSB(x2);
	}
	return false;
}

func ( self *GenMoves ) getTotMoves() uint64 {
	return self.numMoves + self.numMovesq;
}

func ( self *GenMoves ) performRankFileCapture(piece uint, enemies uint64, side uint, allpieces uint64) bool {

	var x2 = self.chessboard[piece];
	for ; x2 != 0; {
		var position = BITScanForward(x2);
		var rankFile = self.getRankFile(position, allpieces) & enemies;
		for ; rankFile != 0; {
			if self.pushmove(STANDARD_MOVE_MASK, position, BITScanForward(rankFile), side, NO_PROMOTION, piece) {
				return true;
			}
			RESET_LSB(rankFile);
		}
		RESET_LSB(x2);
	}
	return false;
}

func ( self *GenMoves ) performPawnCapture(enemies uint64, side uint) bool {
	if self.chessboard[side] == 0 {
		if self.chessboard[ENPASSANT_IDX] != NO_ENPASSANT {
			self.updateZobristKey(13, self.chessboard[ENPASSANT_IDX]);
		}
		self.chessboard[ENPASSANT_IDX] = NO_ENPASSANT;
		return false;
	}
	var GG int
	var x uint64
	if side != 0 {
		x = (self.chessboard[side] << 7) & TABCAPTUREPAWN_LEFT & enemies;
		GG = -7;
	} else {
		x = (self.chessboard[side] >> 7) & TABCAPTUREPAWN_RIGHT & enemies;
		GG = 7;
	};
	for ; x != 0; {
		var o = BITScanForward(x);
		if (side != 0 && o > 55) || (side == 0 && o < 8) {
			//PROMOTION
			if self.pushmove(PROMOTION_MOVE_MASK, o + GG, o, side, (QUEEN_BLACK + side), side) {
				return true; //queen
			}
			if self.perftMode == true {
				if self.pushmove(PROMOTION_MOVE_MASK, o + GG, o, side, (KNIGHT_BLACK + side), side) {
					return true; //knight
				}
				if self.pushmove(PROMOTION_MOVE_MASK, o + GG, o, side, (ROOK_BLACK + side), side) {
					return true; //rock
				}
				if self.pushmove(PROMOTION_MOVE_MASK, o + GG, o, side, (BISHOP_BLACK + side), side) {
					return true; //bishop
				}
			}
		} else if self.pushmove(STANDARD_MOVE_MASK, o + GG, o, side, NO_PROMOTION, side) {
			return true;
		}
		RESET_LSB(x);
	};
	if side != 0 {
		GG = -9;
		x = (self.chessboard[side] << 9) & TABCAPTUREPAWN_RIGHT & enemies;
	} else {
		GG = 9;
		x = (self.chessboard[side] >> 9) & TABCAPTUREPAWN_LEFT & enemies;
	};
	for ; x != 0; {
		var o = BITScanForward(x);
		if (side != 0 && o > 55) || (side == 0 && o < 8) {
			//PROMOTION
			if self.pushmove(PROMOTION_MOVE_MASK, o + GG, o, side, (QUEEN_BLACK + side), side) {
				return true; //queen
			}
			if self.perftMode == true {
				if self.pushmove(PROMOTION_MOVE_MASK, o + GG, o, side, (KNIGHT_BLACK + side), side) {
					return true; //knight
				}
				if self.pushmove(PROMOTION_MOVE_MASK, o + GG, o, side, (BISHOP_BLACK + side), side) {
					return true; //bishop
				}
				if self.pushmove(PROMOTION_MOVE_MASK, o + GG, o, side, (ROOK_BLACK + side), side) {
					return true; //rock
				}
			}
		} else if self.pushmove(STANDARD_MOVE_MASK, o + GG, o, side, NO_PROMOTION, side) {
			return true;
		}
		RESET_LSB(x);
	};
	//ENPASSANT
	if self.chessboard[ENPASSANT_IDX] != NO_ENPASSANT {
		x = ENPASSANT_MASK[side ^ 1][self.chessboard[ENPASSANT_IDX]] & self.chessboard[side];
		for ; x != 0; {
			var o = BITScanForward(x);
			var ff uint64
			if side != 0 {
				ff = self.chessboard[ENPASSANT_IDX] + 8
			} else {
				ff = self.chessboard[ENPASSANT_IDX] - 8
			};
			self.pushmove(ENPASSANT_MOVE_MASK, o, ff, side, NO_PROMOTION, side);
			RESET_LSB(x);
		}
		self.updateZobristKey(13, self.chessboard[ENPASSANT_IDX]);
		self.chessboard[ENPASSANT_IDX] = NO_ENPASSANT;
	}
	return false;
}

func ( self *GenMoves ) performPawnShift(side uint, xallpieces uint64) {
	var tt int
	var x = self.chessboard[side];
	if x & PAWNS_JUMP[side] != 0 {
		self.checkJumpPawn(side, x, xallpieces);
	}
	if side != 0 {
		x <<= 8;
		tt = -8;
	} else {
		tt = 8;
		x >>= 8;
	};
	x &= xallpieces;
	for ; x != 0; {
		var o = BITScanForward(x);

		if o > 55 || o < 8 {
			self.pushmove(PROMOTION_MOVE_MASK, o + tt, o, side, (QUEEN_BLACK + side), side);
			if self.perftMode == true {
				self.pushmove(PROMOTION_MOVE_MASK, o + tt, o, side, (KNIGHT_BLACK + side), side);
				self.pushmove(PROMOTION_MOVE_MASK, o + tt, o, side, (BISHOP_BLACK + side), side);
				self.pushmove(PROMOTION_MOVE_MASK, o + tt, o, side, (ROOK_BLACK + side), side);
			}
		} else {
			self.pushmove(STANDARD_MOVE_MASK, o + tt, o, side, NO_PROMOTION, side);
		}
		RESET_LSB(x);
	};
}

func ( self *GenMoves ) clearKillerHeuristic() {
	//	self.killerHeuristic = [[0; 64]; 64];
	//memset(killerHeuristic, 0, sizeof (killerHeuristic));
}

func ( self *GenMoves ) performDiagShift(piece uint, side uint, allpieces uint64) {

	var x2 = self.chessboard[piece];
	for ; x2 != 0; {
		var position = BITScanForward(x2);
		var diag = self.getDiagonalAntiDiagonal(position, allpieces) & !allpieces;
		for ; diag != 0; {
			self.pushmove(STANDARD_MOVE_MASK, position, BITScanForward(diag), side, NO_PROMOTION, piece);
			RESET_LSB(diag);
		}
		RESET_LSB(x2);
	}
}

func ( self *GenMoves ) performRankFileShift(piece uint, side uint, allpieces uint64) {

	var x2 = self.chessboard[piece];
	for ; x2 != 0; {
		var position = BITScanForward(x2);
		var rankFile = self.getRankFile(position, allpieces) & !allpieces;
		for ; rankFile != 0; {
			self.pushmove(STANDARD_MOVE_MASK, position, BITScanForward(rankFile), side, NO_PROMOTION, piece);
			RESET_LSB(rankFile);
		}
		RESET_LSB(x2);
	}
}

func ( self *GenMoves ) makemove(mov *_Tmove, rep bool, checkInCheck bool) bool {
	var pieceFrom uint = SQUARE_FREE;
	var posTo uint8;
	var posFrom uint8;
	var movecapture = SQUARE_FREE;
	var rightCastleOld = self.chessboard[RIGHT_CASTLE_IDX];
	if mov.typee & 0xc == 0 {
		//no castle
		posTo = mov.to;
		posFrom = mov.from;
		movecapture = mov.capturedPiece;

		pieceFrom = mov.pieceFrom;
		if (mov.typee & 0x3) == PROMOTION_MOVE_MASK {
			self.chessboard[pieceFrom] &= NOTPOW2[posFrom];
			self.updateZobristKey(pieceFrom, posFrom);

			self.chessboard[mov.promotionPiece] |= POW2[posTo];
			self.updateZobristKey(mov.promotionPiece, posTo);
		} else {
			self.chessboard[pieceFrom] = (self.chessboard[pieceFrom] | POW2[posTo]) & NOTPOW2[posFrom];
			self.updateZobristKey(pieceFrom, posFrom);
			self.updateZobristKey(pieceFrom, posTo);
		}
		if movecapture != SQUARE_FREE {
			if (mov.typee & 0x3) != ENPASSANT_MOVE_MASK {
				self.chessboard[movecapture] &= NOTPOW2[posTo];
				self.updateZobristKey(movecapture, posTo);
			} else {
				//en passant

				if mov.side != 0 {
					self.chessboard[movecapture] &= NOTPOW2[posTo - 8];
					self.updateZobristKey(movecapture, posTo - 8);
				} else {
					self.chessboard[movecapture] &= NOTPOW2[posTo + 8];
					self.updateZobristKey(movecapture, posTo + 8);
				}
			}
		}
		//lost castle right
		switch pieceFrom{
		case KING_WHITE :{
			self.chessboard[RIGHT_CASTLE_IDX] &= 0xcf;
		}

		case KING_BLACK :{
			self.chessboard[RIGHT_CASTLE_IDX] &= 0x3f;
		}
		case ROOK_WHITE :{
			if posFrom == 0 {
				self.chessboard[RIGHT_CASTLE_IDX] &= 0xef;
			} else if posFrom == 7 {
				self.chessboard[RIGHT_CASTLE_IDX] &= 0xdf;
			}
		}
		case ROOK_BLACK :{
			if posFrom == 56 {
				self.chessboard[RIGHT_CASTLE_IDX] &= 0xbf;
			} else if posFrom == 63 {
				self.chessboard[RIGHT_CASTLE_IDX] &= 0x7f;
			}
		}
		//en passant

		case PAWN_WHITE :{
			if (RANK_1 & POW2[posFrom]) != 0 && (RANK_3 & POW2[posTo]) != 0 {
				self.chessboard[ENPASSANT_IDX] = posTo;
				self.updateZobristKey(13, self.chessboard[ENPASSANT_IDX]);
			}
		}

		case PAWN_BLACK :{
			if (RANK_6 & POW2[posFrom]) != 0 && (RANK_4 & POW2[posTo]) != 0 {
				self.chessboard[ENPASSANT_IDX] = posTo;
				self.updateZobristKey(13, self.chessboard[ENPASSANT_IDX]);
			}
		}

		}
	} else {
		//castle
		self.performCastle(mov.side, mov.typee);
		if mov.side == WHITE {
			self.chessboard[RIGHT_CASTLE_IDX] &= 0xcf;
		} else {
			self.chessboard[RIGHT_CASTLE_IDX] &= 0x3f;
		}
	}
	var x2 = rightCastleOld ^ self.chessboard[RIGHT_CASTLE_IDX];
	for ; x2 != 0; {
		var position = BITScanForward(x2);
		self.updateZobristKey(14, position);
		RESET_LSB(x2);
	}
	//if rep == true {
	//	if movecapture != SQUARE_FREE || pieceFrom == WHITE || pieceFrom == BLACK || mov.typee & 0xcuint64 != 0 {
	//		self.pushStackMove1(0);
	//	}
	//	self.pushStackMove1(self.chessboard[ZOBRISTKEY_IDX]);
	//}
	if (self.forceCheck || (checkInCheck == true && self.perftMode == false)) && ((mov.side == WHITE && self.inCheck1(WHITE)) || (mov.side == BLACK && self.inCheck1(BLACK))) {
		return false;
	}
	return true;
}

func ( self *GenMoves ) incListId() {
	self.listId = self.listId + 1;
	//#ifdef DEBUG_MODE
	//if (listId < 0 || listId >= MAX_PLY) {
	//display();
	//}
	//debug_assert!_RANGE(listId, 0, MAX_PLY - 1);
	//#endif
}

func ( self *GenMoves ) display() {
	self.display();
}

func ( self *GenMoves ) decListId() {

	self.gen_list[self.listId].size = 0;
	self.listId = self.listId - 1;
}
//
func ( self *GenMoves ) getListSize() int {
	return self.gen_list[self.listId].size;
}
//
//func ( self *GenMoves ) pushStackMove() {
//	self.pushStackMove1(self.chessboard[ZOBRISTKEY_IDX]);
//}

func ( self *GenMoves ) resetList() {
	self.gen_list[self.listId].size = 0;
}
//
//func ( self *GenMoves ) incKillerHeuristic(from uint, to uint, value int) {
//	if self.getRunning() == 0 {
//		return;
//	}
//
//
//	//        debug_assert!(self.killerHeuristic[from][to] <= self.killerHeuristic[from][to] + value);
//	//self.killerHeuristic[from][to] += value;
//}

func ( self *GenMoves ) getNextMove(list _TmoveP) *_Tmove {
	var gen_list1 = &list.moveList;

	var listcount uint = list.size;
	var bestId = -1;

	var bestScore int;
	var j uint;
	for j := 0; j < listcount; j++ {
		if gen_list1[j].used == false {
			bestId = j;
			bestScore = gen_list1[bestId].score;
			break;
		}
	}
	if bestId == -1 {
		return nil;
	}
	for i := j + 1; i < listcount; i++ {
		if gen_list1[i].used == true && gen_list1[i].score > bestScore {
			bestId = i;
			bestScore = gen_list1[bestId].score;
		}
	}
	gen_list1[bestId].used = true;
	return &gen_list1[bestId];
}

func ( self *GenMoves ) isAttacked(side uint, position uint, allpieces uint64) uint64 {
	return self.getAttackers(side, true, position, allpieces)
}

func ( self *GenMoves ) getAllAttackers(side uint, position uint, allpieces uint64) uint64 {
	return self.getAttackers(side, false, position, allpieces)
}

func ( self *GenMoves ) getMobilityRook(position uint, enemies uint64, friends uint64) uint {

	return self.performRankFileCaptureAndShiftCount(position, enemies, enemies | friends)
}

func ( self *GenMoves ) getMobilityPawns(side uint, ep uint64, ped_friends uint64, enemies uint64, xallpieces uint64) uint {

	if ep == NO_ENPASSANT {
		return 0
	}
	if bitCount((ENPASSANT_MASK[side ^ 1][ep] & self.chessboard[side])) + side == WHITE {
		return bitCount(((ped_friends << 8) & xallpieces)) + bitCount(((((ped_friends & TABJUMPPAWN) << 8) & xallpieces) << 8) & xallpieces) + bitCount((self.chessboard[side] << 7) & TABCAPTUREPAWN_LEFT & enemies) + bitCount((self.chessboard[side] << 9) & TABCAPTUREPAWN_RIGHT & enemies)
	} else {
		return bitCount(((ped_friends >> 8) & xallpieces)) + bitCount(((((ped_friends & TABJUMPPAWN) >> 8) & xallpieces) >> 8) & xallpieces) + bitCount((self.chessboard[side] >> 7) & TABCAPTUREPAWN_RIGHT & enemies) + bitCount((self.chessboard[side] >> 9) & TABCAPTUREPAWN_LEFT & enemies);
	}
}

func ( self *GenMoves ) getMobilityCastle(side uint, allpieces uint64) int {

	var count = 0;
	if side == WHITE {
		if POW2_3 & self.chessboard[KING_WHITE] != 0 && allpieces & 0x6 == 0 && self.chessboard[RIGHT_CASTLE_IDX] & RIGHT_KING_CASTLE_WHITE_MASK != 0 && self.chessboard[ROOK_WHITE] & POW2_0 != 0 && self.isAttacked(WHITE, 1, allpieces) == 0 && self.isAttacked(WHITE, 2, allpieces) == 0 && self.isAttacked(WHITE, 3, allpieces) == 0 {
			count = count + 1;
		}
		if POW2_3 & self.chessboard[KING_WHITE] != 0 && !(allpieces & 0x70) == 0 && self.chessboard[RIGHT_CASTLE_IDX] & RIGHT_QUEEN_CASTLE_WHITE_MASK != 0 && self.chessboard[ROOK_WHITE] & POW2_7 != 0 && self.isAttacked(WHITE, 3, allpieces) == 0 && self.isAttacked(WHITE, 4, allpieces) == 0 && self.isAttacked(WHITE, 5, allpieces) == 0 {
			count = count + 1;
		}
	} else {
		if POW2_59 & self.chessboard[KING_BLACK] != 0 && self.chessboard[RIGHT_CASTLE_IDX] & RIGHT_KING_CASTLE_BLACK_MASK != 0 && (allpieces & 0x600000000000000) == 0 && self.chessboard[ROOK_BLACK] & POW2_56 != 0 && self.isAttacked(BLACK, 57, allpieces) == 0 && self.isAttacked(BLACK, 58, allpieces) == 0 && self.isAttacked(BLACK, 59, allpieces) == 0 {
			count = count + 1;
		}
		if POW2_59 & self.chessboard[KING_BLACK] != 0 && self.chessboard[RIGHT_CASTLE_IDX] & RIGHT_QUEEN_CASTLE_BLACK_MASK != 0 && (allpieces & 0x7000000000000000) == 0 && self.chessboard[ROOK_BLACK] & POW2_63 != 0 && self.isAttacked(BLACK, 59, allpieces) == 0 && self.isAttacked(BLACK, 60, allpieces) == 0 && self.isAttacked(BLACK, 61, allpieces) == 0 {
			count = count + 1;
		}
	}
	return count;
}

func ( self *GenMoves ) getMobilityQueen(position uint, enemies uint64, allpieces uint64) uint {

	return self.performRankFileCaptureAndShiftCount(position, enemies, allpieces) +
		bitCount(self.getDiagShiftAndCapture(position, enemies, allpieces))
}

func ( self *GenMoves ) inCheck(side uint, typee uint, from uint, to uint, pieceFrom uint, pieceTo uint, promotionPiece uint) uint64 {
	//#ifdef DEBUG_MODE
	//_Tchessboard a;
	//memcpy( & a, self.chessboard, sizeof (_Tchessboard));
	//#endif

	var result uint64;
	var g3 uint64 = typee & 0x3;
	switch  g3{
	case
		STANDARD_MOVE_MASK :{
		var from1 = self.chessboard[pieceFrom];
		var to1 uint64 = -1;


		if pieceTo != SQUARE_FREE {
			to1 = self.chessboard[pieceTo];
			self.chessboard[pieceTo] &= NOTPOW2[to];
		}
		self.chessboard[pieceFrom] &= NOTPOW2[from];
		self.chessboard[pieceFrom] |= POW2[to];

		result = self.isAttacked(side, BITScanForward(self.chessboard[KING_BLACK + side]), self.getBitmap(BLACK) | self.getBitmap(WHITE));
		self.chessboard[pieceFrom] = from1;
		if pieceTo != SQUARE_FREE {
			self.chessboard[pieceTo] = to1;
		}
	}

	case PROMOTION_MOVE_MASK :{
		var to1 uint64 = 0;
		if pieceTo != SQUARE_FREE {
			to1 = self.chessboard[pieceTo];
		}
		var from1 uint64 = self.chessboard[pieceFrom];
		var p1 = self.chessboard[promotionPiece];
		self.chessboard[pieceFrom] &= NOTPOW2[from];
		if pieceTo != SQUARE_FREE {
			self.chessboard[pieceTo] &= NOTPOW2[to];
		}
		self.chessboard[promotionPiece] = self.chessboard[promotionPiece] | POW2[to];
		result = self.isAttacked(side, BITScanForward(self.chessboard[KING_BLACK + side]), self.getBitmap(BLACK) | self.getBitmap(WHITE));
		if pieceTo != SQUARE_FREE {
			self.chessboard[pieceTo] = to1;
		}
		self.chessboard[pieceFrom] = from1;
		self.chessboard[promotionPiece] = p1;
	}

	case ENPASSANT_MOVE_MASK :{
		var to1 = self.chessboard[side ^ 1];
		var from1 = self.chessboard[side];
		self.chessboard[side] &= NOTPOW2[from];
		self.chessboard[side] |= POW2[to];
		if side != 0 {
			self.chessboard[side ^ 1] &= NOTPOW2[to - 8];
		} else {
			self.chessboard[side ^ 1] &= NOTPOW2[to + 8];
		}
		result = self.isAttacked(side, BITScanForward(self.chessboard[KING_BLACK + side]), self.getBitmap(BLACK) | self.getBitmap(WHITE));
		self.chessboard[side ^ 1] = to1;
		self.chessboard[side] = from1;
	}

	}

	//    #ifdef DEBUG_MODE

	//    #endif
	return result;
}

func ( self *GenMoves ) performCastle(side uint, typee uint64) {

	if side == WHITE {
		if typee & KING_SIDE_CASTLE_MOVE_MASK != 0 {

			self.updateZobristKey(KING_WHITE, 3);
			self.updateZobristKey(KING_WHITE, 1);
			self.chessboard[KING_WHITE] = self.chessboard[KING_WHITE] | POW2_1 & NOTPOW2_3;
			self.updateZobristKey(ROOK_WHITE, 2);
			self.updateZobristKey(ROOK_WHITE, 0);
			self.chessboard[ROOK_WHITE] = self.chessboard[ROOK_WHITE] | POW2_2 & NOTPOW2_0;
		} else {

			self.chessboard[KING_WHITE] = self.chessboard[KING_WHITE] | POW2_5 & NOTPOW2_3;
			self.updateZobristKey(KING_WHITE, 5);
			self.updateZobristKey(KING_WHITE, 3);
			self.chessboard[ROOK_WHITE] = self.chessboard[ROOK_WHITE] | POW2_4 & NOTPOW2_7;
			self.updateZobristKey(ROOK_WHITE, 4);
			self.updateZobristKey(ROOK_WHITE, 7);
		}
	} else {
		if typee & KING_SIDE_CASTLE_MOVE_MASK != 0 {

			self.chessboard[KING_BLACK] = self.chessboard[KING_BLACK] | POW2_57 & NOTPOW2_59;
			self.updateZobristKey(KING_BLACK, 57);
			self.updateZobristKey(KING_BLACK, 59);
			self.chessboard[ROOK_BLACK] = self.chessboard[ROOK_BLACK] | POW2_58 & NOTPOW2_56;
			self.updateZobristKey(ROOK_BLACK, 58);
			self.updateZobristKey(ROOK_BLACK, 56);
		} else {

			self.chessboard[KING_BLACK] = self.chessboard[KING_BLACK] | POW2_61 & NOTPOW2_59;
			self.updateZobristKey(KING_BLACK, 61);
			self.updateZobristKey(KING_BLACK, 59);
			self.chessboard[ROOK_BLACK] = self.chessboard[ROOK_BLACK] | POW2_60 & NOTPOW2_63;
			self.updateZobristKey(ROOK_BLACK, 60);
			self.updateZobristKey(ROOK_BLACK, 63);
		}
	}
}

func ( self *GenMoves ) unPerformCastle(side uint, typee uint64) {
	if side == WHITE {
		if typee & KING_SIDE_CASTLE_MOVE_MASK != 0 {

			self.chessboard[KING_WHITE] = (self.chessboard[KING_WHITE] | POW2_3) & NOTPOW2_1;
			self.chessboard[ROOK_WHITE] = (self.chessboard[ROOK_WHITE] | POW2_0) & NOTPOW2_2;
		} else {
			self.chessboard[KING_WHITE] = (self.chessboard[KING_WHITE] | POW2_3) & NOTPOW2_5;
			self.chessboard[ROOK_WHITE] = (self.chessboard[ROOK_WHITE] | POW2_7) & NOTPOW2_4;
		}
	} else {
		if typee & KING_SIDE_CASTLE_MOVE_MASK != 0 {
			self.chessboard[KING_BLACK] = (self.chessboard[KING_BLACK] | POW2_59) & NOTPOW2_57;
			self.chessboard[ROOK_BLACK] = (self.chessboard[ROOK_BLACK] | POW2_56) & NOTPOW2_58;
		} else {
			self.chessboard[KING_BLACK] = (self.chessboard[KING_BLACK] | POW2_59) & NOTPOW2_61;
			self.chessboard[ROOK_BLACK] = (self.chessboard[ROOK_BLACK] | POW2_63) & NOTPOW2_60;
		}
	}
}

func ( self *GenMoves ) tryAllCastle(side uint, allpieces uint64) {

	if side == WHITE {
		if POW2_3 & self.chessboard[KING_WHITE] != 0 && allpieces & 0x6 == 0 && self.chessboard[RIGHT_CASTLE_IDX] & RIGHT_KING_CASTLE_WHITE_MASK != 0  && self.chessboard[ROOK_WHITE] & POW2_0 != 0  && 0 == self.isAttacked(WHITE, 1, allpieces) && 0 == self.isAttacked(WHITE, 2, allpieces) && 0 == self.isAttacked(WHITE, 3, allpieces) {
			self.pushmove(KING_SIDE_CASTLE_MOVE_MASK, -1, -1, WHITE, NO_PROMOTION, -1);
		}
		if POW2_3 & self.chessboard[KING_WHITE] != 0 && allpieces & 0x70 == 0 && self.chessboard[RIGHT_CASTLE_IDX] & RIGHT_QUEEN_CASTLE_WHITE_MASK != 0 && self.chessboard[ROOK_WHITE] & POW2_7 != 0  && 0 == self.isAttacked(WHITE, 3, allpieces) && 0 == self.isAttacked(WHITE, 4, allpieces) && 0 == self.isAttacked(WHITE, 5, allpieces) {
			self.pushmove(QUEEN_SIDE_CASTLE_MOVE_MASK, -1, -1, WHITE, NO_PROMOTION, -1);
		}
	} else {
		if POW2_59 & self.chessboard[KING_BLACK] != 0 && self.chessboard[RIGHT_CASTLE_IDX] & RIGHT_KING_CASTLE_BLACK_MASK != 0 && 0 == (allpieces & 0x600000000000000) && self.chessboard[ROOK_BLACK] & POW2_56 != 0 && 0 == self.isAttacked(BLACK, 57, allpieces) && 0 == self.isAttacked(BLACK, 58, allpieces) && 0 == self.isAttacked(BLACK, 59, allpieces) {
			self.pushmove(KING_SIDE_CASTLE_MOVE_MASK, -1, -1, BLACK, NO_PROMOTION, -1);
		}
		if POW2_59 & self.chessboard[KING_BLACK] != 0&& self.chessboard[RIGHT_CASTLE_IDX] & RIGHT_QUEEN_CASTLE_BLACK_MASK != 0 && 0 == (allpieces & 0x7000000000000000) && self.chessboard[ROOK_BLACK] & POW2_63 != 0 && 0 == self.isAttacked(BLACK, 59, allpieces) && 0 == self.isAttacked(BLACK, 60, allpieces) && 0 == self.isAttacked(BLACK, 61, allpieces) {
			self.pushmove(QUEEN_SIDE_CASTLE_MOVE_MASK, -1, -1, BLACK, NO_PROMOTION, -1);
		}
	}
}

func ( self *GenMoves ) pushmove(typee uint64, from uint, to uint, side uint, promotionPiece int, pieceFrom uint) bool {

	var piece_captured = SQUARE_FREE;
	var res = false;
	if ((typee & 0x3) != ENPASSANT_MOVE_MASK) && 0 == (typee & 0xc) {
		if side == WHITE {
			piece_captured = self.getPieceAt(BLACK, POW2[to])
		} else {
			piece_captured = self.getPieceAt(WHITE, POW2[to])
		}

		if piece_captured == KING_BLACK + (side ^ 1) {
			res = true;
		}
	} else if typee & 0xc == 0 {
		//no castle
		piece_captured = side ^ 1;
	}
	if (typee & 0xc) != 0 && (self.forceCheck || self.perftMode) {
		//no castle
		if side == WHITE && self.inCheck(WHITE, typee, from, to, pieceFrom, piece_captured, promotionPiece) != 0 {
			return false;
		}
		if side == BLACK && self.inCheck(BLACK, typee, from, to, pieceFrom, piece_captured, promotionPiece) != 0 {
			return false;
		}
	}
	var mos*_Tmove;


	mos = &self.gen_list[self.listId].moveList[self.getListSize()];
	self.gen_list[self.listId].size = self.gen_list[self.listId].size + 1;
	mos.typee = self.chessboard[RIGHT_CASTLE_IDX] | typee;
	mos.side = side;
	mos.capturedPiece = piece_captured;
	if typee & 0x3 != 0 {
		mos.from = from;
		mos.to = to;
		mos.pieceFrom = pieceFrom;
		mos.promotionPiece = promotionPiece;
		if self.perftMode == false {
			if res == true {
				mos.score = _INFINITE;
			} else {

				//mos.score = self.killerHeuristic[from][to];
				if PIECES_VALUE[piece_captured] >= PIECES_VALUE[pieceFrom] {
					mos.score = mos.score + (PIECES_VALUE[piece_captured] - PIECES_VALUE[pieceFrom]) * 2;
				} else {
					mos.score = mos.score + PIECES_VALUE[piece_captured];
				}
			}
		}
	} else if typee & 0xc != 0 {
		//castle

		mos.score = 100;
	}
	mos.used = false;

	return res;
}

func ( self *GenMoves ) getMove(i uint) *_Tmove {
	return self.gen_list[self.listId].moveList[i]
}
//
//func ( self *GenMoves ) setRunning(t int) {
//	self.running = t;
//}
//
//func ( self *GenMoves ) getRunning() int {
//	return self.running;
//}

func ( self *GenMoves ) inCheck1(side uint) bool {
	return self.isAttacked(side, BITScanForward(self.chessboard[KING_BLACK + side]), self.getBitmap(BLACK) | self.getBitmap(WHITE));
}

//func ( self *GenMoves ) setKillerHeuristic(from uint, to uint, value int) {
//	if self.getRunning() != 0 {
//		self.killerHeuristic[from][to] = value;
//	}
//}

///////////////////////////////////// private:


func ( self *GenMoves ) checkJumpPawn(side uint, xx uint64, xallpieces uint64) {
	var x = xx;
	x &= TABJUMPPAWN;
	if side != 0 {
		x = (((x << 8) & xallpieces) << 8) & xallpieces;
	} else {
		x = (((x >> 8) & xallpieces) >> 8) & xallpieces;
	};
	for ; x != 0; {
		var o = BITScanForward(x);
		var rr int;
		if side != 0 {
			rr = -16;
		} else {
			rr = 16;
		}

		self.pushmove(STANDARD_MOVE_MASK, ((o + rr)), o, side, NO_PROMOTION, side);
		RESET_LSB(x);
	}
}

func ( self *GenMoves ) performRankFileCaptureAndShiftCount(position uint, enemies uint64, allpieces uint64) uint {

	var rankFile uint64 = self.getRankFile(position, allpieces);
	rankFile = (rankFile & enemies) | (rankFile & !allpieces);
	return bitCount(rankFile)
}
/*
func ( self *GenMoves ) popStackMove() {

	self.repetitionMapCount = self.repetitionMapCount - 1;
	if self.repetitionMapCount != 0 && self.repetitionMap[self.repetitionMapCount - 1] == 0 {
		self.repetitionMapCount = self.repetitionMapCount - 1;
	}
}

func ( self *GenMoves ) pushStackMove1(key uint64) {
	self.repetitionMap[self.repetitionMapCount] = key;
	self.repetitionMapCount = self.repetitionMapCount + 1;
}*/

func ( self *GenMoves ) getAttackers(side uint, exitOnFirst bool, position uint, allpieces uint64) uint64 {

	//knight
	var attackers = KNIGHT_MASK[position] & self.chessboard[KNIGHT_BLACK + (side ^ 1)];
	if exitOnFirst == true && attackers != 0 {
		return 1
	};
	//king
	attackers |= NEAR_MASK1[position] & self.chessboard[KING_BLACK + (side ^ 1)];
	if exitOnFirst == true && attackers != 0 {
		return 1
	};
	//pawn
	attackers |= PAWN_FORK_MASK[side][position] & self.chessboard[PAWN_BLACK + (side ^ 1)];
	if exitOnFirst == true && attackers != 0 {
		return 1
	};
	//bishop queen
	var enemies = self.chessboard[BISHOP_BLACK + (side ^ 1)] | self.chessboard[QUEEN_BLACK + (side ^ 1)];
	var nuovo = self.getDiagonalAntiDiagonal(position, allpieces) & enemies;
	for ; nuovo != 0; {
		var bound = BITScanForward(nuovo);
		attackers |= POW2[bound];
		if exitOnFirst == true && attackers != 0 {
			return 1
		};
		RESET_LSB(nuovo);
	}
	enemies = self.chessboard[ROOK_BLACK + (side ^ 1)] | self.chessboard[QUEEN_BLACK + (side ^ 1)];
	nuovo = self.getRankFile(position, allpieces) & enemies;
	for ; nuovo != 0; {
		var bound = BITScanForward(nuovo);
		attackers |= POW2[bound];
		if exitOnFirst == true && attackers != 0 {
			return 1
		};
		RESET_LSB(nuovo);
	}
	return attackers;
}
