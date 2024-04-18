package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Position", func() {
	Describe("::String", func() {
		It("represents the init board clearly", func() {
			pos := InitPos()
			expStr := "" +
				"w KQkq - 0 1\n" +
				"8 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ \n" +
				"7 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙ \n" +
				"6   ░░  ░░  ░░  ░░\n" +
				"5 ░░  ░░  ░░  ░░  \n" +
				"4   ░░  ░░  ░░  ░░\n" +
				"3 ░░  ░░  ░░  ░░  \n" +
				"2 ♟︎ ♟︎ ♟︎ ♟︎ ♟︎ ♟︎ ♟︎ ♟︎ \n" +
				"1 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ \n" +
				"  1 2 3 4 5 6 7 8 "
			Expect(pos.String()).To(Equal(expStr))
		})
	})
	Describe("#FromFEN", func() {
		When("the FEN is invalid", func() {
			When("the FEN does not have enough segments", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0"
					Expect(FromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("the FEN has too many segments", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 1"
					Expect(FromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("the FEN pieces contain too many (9) rows", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR/8 w KQkq - 0 1"
					Expect(FromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("the FEN pieces contain too few (7) rows", func() {
				It("returns an error", func() {
					fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP w KQkq - 0 1"
					Expect(FromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("a row in the FEN pieces contain too many (9) files", func() {
				It("returns an error", func() {
					fen := "rnbqkbnrr/ppppppppp/9/9/9/9/PPPPPPPPP/RNBQKBNRR w KQkq - 0 1"
					Expect(FromFEN(fen)).Error().To(HaveOccurred())
				})
			})
			When("a row in the FEN pieces contain too few (7) files", func() {
				It("returns an error", func() {
					fen := "rnbqkbn/ppppppp/7/7/7/7/PPPPPPP/RNBQKBN w KQkq - 0 1"
					Expect(FromFEN(fen)).Error().To(HaveOccurred())
				})
			})
		})
		When("the FEN is valid", func() {
			It("returns a Position with the pieces as described in the FEN", func() {
				fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
				pos, posErr := FromFEN(fen)
				Expect(posErr).To(Succeed())
				Expect(*pos).To(Equal(*InitPos()))
			})
		})
	})
	Describe("::ToFEN", func() {
		It("serializes the position into the FEN format", func() {
			fen := InitPos().FEN()
			Expect(fen).To(Equal("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"))
		})
	})
	Describe("::MakeMove", func() {
		var pos *Position
		var prevHash ZHash
		BeforeEach(func() {
			var posErr error
			pos, posErr = FromFEN("3k4/1b6/8/8/8/8/3P2B1/4K2R w K - 0 1")
			Expect(posErr).ToNot(HaveOccurred())
			expPieceBBs := [N_PIECES]Bitboard{
				0b11110111_11111101_11111111_11111111_11111111_11111111_10110111_01101111,
				BBWithSquares(SQ_D2),
				0,
				BBWithSquares(SQ_G2),
				BBWithSquares(SQ_H1),
				0,
				BBWithSquares(SQ_E1),
				0,
				0,
				BBWithSquares(SQ_B7),
				0,
				0,
				BBWithSquares(SQ_D8),
			}
			Expect(pos.pieceBitboards).To(Equal(expPieceBBs))
			expPieces := [N_SQUARES]Piece{}
			expPieces[SQ_E1] = W_KING
			expPieces[SQ_H1] = W_ROOK
			expPieces[SQ_D2] = W_PAWN
			expPieces[SQ_G2] = W_BISHOP
			expPieces[SQ_B7] = B_BISHOP
			expPieces[SQ_D8] = B_KING
			Expect(pos.pieces).To(Equal(expPieces))
			Expect(pos.colorBitboards[WHITE]).To(Equal(BBWithSquares(SQ_E1, SQ_H1, SQ_D2, SQ_G2)))
			Expect(pos.colorBitboards[BLACK]).To(Equal(BBWithSquares(SQ_B7, SQ_D8)))
			Expect(pos.material[B_BISHOP]).To(BeEquivalentTo(1))
			Expect(pos.hash).ToNot(BeEquivalentTo(0))
			prevHash = pos.hash
		})
		When("the move is a pawn move", func() {
			It("updates the position pieces", func() {
				pos.MakeMove(NewNormalMove(SQ_D2, SQ_D4))
				Expect(pos.pieces[SQ_D2]).To(Equal(EMPTY))
				Expect(pos.pieces[SQ_D4]).To(Equal(W_PAWN))
				expEmptyBB := Bitboard(0b11110111_11111101_11111111_11111111_11110111_11111111_10111111_01101111)
				Expect(pos.pieceBitboards[EMPTY]).To(Equal(expEmptyBB))
				Expect(pos.pieceBitboards[W_PAWN]).To(Equal(BBWithSquares(SQ_D4)))
				Expect(pos.colorBitboards[WHITE]).To(Equal(BBWithSquares(SQ_D4, SQ_G2, SQ_E1, SQ_H1)))
				Expect(pos.colorBitboards[BLACK]).To(Equal(BBWithSquares(SQ_D8, SQ_B7)))
				Expect(pos.isWhiteTurn).To(BeFalse())
			})
			It("updates the hash", func() {
				pos.MakeMove(NewNormalMove(SQ_D2, SQ_D4))
				Expect(pos.hash).ToNot(Equal(prevHash))
			})
			It("updates the ply", func() {
				pos.MakeMove(NewNormalMove(SQ_D2, SQ_D4))
				Expect(pos.ply).To(BeEquivalentTo(1))
			})
			It("updates a copy of frozenPos", func() {
				oldFP := pos.frozenPos
				pos.MakeMove(NewNormalMove(SQ_D2, SQ_D4))
				fp := pos.frozenPos
				Expect(fp.Rule50).To(BeEquivalentTo(0))
				Expect(fp.EnPassantSq).To(Equal(SQ_D3))
				expCastleRights := [N_CASTLE_RIGHTS]bool{true, false, false, false}
				Expect(fp.CastleRights).To(Equal(expCastleRights))
				Expect(oldFP).ToNot(Equal(fp))
			})
			It("adds the new position to the repetitions map", func() {
				pos.MakeMove(NewNormalMove(SQ_D2, SQ_D4))
				Expect(pos.repetitions[pos.hash]).To(BeEquivalentTo(1))
			})
		})
		When("the move is castles", func() {
			It("updates the position pieces", func() {
				pos.MakeMove(NewMove(SQ_E1, SQ_G1, NULL_SQ, EMPTY_PIECE_TYPE, true))
				Expect(pos.pieces[SQ_E1]).To(Equal(EMPTY))
				Expect(pos.pieces[SQ_F1]).To(Equal(W_ROOK))
				Expect(pos.pieces[SQ_G1]).To(Equal(W_KING))
				expEmptyBB := Bitboard(0b11110111_11111101_11111111_11111111_11111111_11111111_10110111_10011111)
				Expect(pos.pieceBitboards[EMPTY]).To(Equal(expEmptyBB))
				Expect(pos.pieceBitboards[W_ROOK]).To(Equal(BBWithSquares(SQ_F1)))
				Expect(pos.pieceBitboards[W_KING]).To(Equal(BBWithSquares(SQ_G1)))
				Expect(pos.colorBitboards[WHITE]).To(Equal(BBWithSquares(SQ_F1, SQ_G1, SQ_D2, SQ_G2)))
				Expect(pos.colorBitboards[BLACK]).To(Equal(BBWithSquares(SQ_B7, SQ_D8)))
				Expect(pos.isWhiteTurn).To(BeFalse())
			})
			It("updates the hash", func() {
				pos.MakeMove(NewMove(SQ_E1, SQ_G1, NULL_SQ, EMPTY_PIECE_TYPE, true))
				Expect(pos.hash).ToNot(Equal(prevHash))
			})
			It("updates a copy of frozenPos", func() {
				oldFP := pos.frozenPos
				pos.MakeMove(NewMove(SQ_E1, SQ_G1, NULL_SQ, EMPTY_PIECE_TYPE, true))
				fp := pos.frozenPos
				Expect(fp).ToNot(Equal(oldFP))
				Expect(fp.Rule50).To(BeEquivalentTo(0))
				expCastleRights := [N_CASTLE_RIGHTS]bool{false, false, false, false}
				Expect(fp.CastleRights).To(Equal(expCastleRights))
				Expect(fp.EnPassantSq).To(Equal(NULL_SQ))
			})
		})
		When("the move is a capture", func() {
			It("updates the position pieces", func() {
				pos.MakeMove(NewNormalMove(SQ_G2, SQ_B7))
				Expect(pos.pieces[SQ_G2]).To(Equal(EMPTY))
				Expect(pos.pieces[SQ_B7]).To(Equal(W_BISHOP))
				Expect(pos.pieceBitboards[W_BISHOP]).To(Equal(BBWithSquares(SQ_B7)))
				Expect(pos.pieceBitboards[B_BISHOP]).To(Equal(Bitboard(0)))
				Expect(pos.colorBitboards[WHITE]).To(Equal(BBWithSquares(SQ_E1, SQ_H1, SQ_D2, SQ_B7)))
				Expect(pos.colorBitboards[BLACK]).To(Equal(BBWithSquares(SQ_D8)))
				Expect(pos.material[B_BISHOP]).To(BeEquivalentTo(0))
				Expect(pos.isWhiteTurn).To(BeFalse())
			})
			It("updates the hash", func() {
				pos.MakeMove(NewNormalMove(SQ_G2, SQ_B7))
				Expect(pos.hash).ToNot(Equal(prevHash))
			})
			It("returns the captured piece", func() {
				capturedPiece := pos.MakeMove(NewNormalMove(SQ_G2, SQ_B7))
				Expect(capturedPiece).To(Equal(B_BISHOP))
			})
			It("updates a copy of frozenPos", func() {
				pos.MakeMove(NewNormalMove(SQ_G2, SQ_B7))
				fp := pos.frozenPos
				Expect(fp.Rule50).To(BeEquivalentTo(0))
				expCastleRights := [N_CASTLE_RIGHTS]bool{true, false, false, false}
				Expect(fp.CastleRights).To(Equal(expCastleRights))
				Expect(fp.EnPassantSq).To(Equal(NULL_SQ))
			})
		})
		When("the move is en passant", func() {
			BeforeEach(func() {
				var posErr error
				pos, posErr = FromFEN("3k4/1b6/8/5Pp1/8/8/3P2B1/4K2R w K g6 1 1")
				Expect(posErr).ToNot(HaveOccurred())
				expPiecesBB := [N_PIECES]Bitboard{
					0b11110111_11111101_11111111_10011111_11111111_11111111_10110111_01101111,
					BBWithSquares(SQ_F5, SQ_D2),
					0,
					BBWithSquares(SQ_G2),
					BBWithSquares(SQ_H1),
					0,
					BBWithSquares(SQ_E1),
					BBWithSquares(SQ_G5),
					0,
					BBWithSquares(SQ_B7),
					0,
					0,
					BBWithSquares(SQ_D8),
				}
				Expect(pos.pieceBitboards).To(Equal(expPiecesBB))
			})
			It("updates the position pieces", func() {
				pos.MakeMove(NewEnPassantMove(SQ_F5, SQ_G6))
				Expect(pos.pieces[SQ_F5]).To(Equal(EMPTY))
				Expect(pos.pieces[SQ_G6]).To(Equal(W_PAWN))
				Expect(pos.pieces[SQ_G5]).To(Equal(EMPTY))
				Expect(pos.pieceBitboards[EMPTY]).To(Equal(Bitboard(0b11110111_11111101_10111111_11111111_11111111_11111111_10110111_01101111)))
				Expect(pos.pieceBitboards[W_PAWN]).To(Equal(BBWithSquares(SQ_D2, SQ_G6)))
				Expect(pos.pieceBitboards[B_PAWN]).To(Equal(Bitboard(0)))
				Expect(pos.colorBitboards[WHITE]).To(Equal(BBWithSquares(SQ_E1, SQ_H1, SQ_D2, SQ_G2, SQ_G6)))
				Expect(pos.colorBitboards[BLACK]).To(Equal(BBWithSquares(SQ_B7, SQ_D8)))
			})
			It("updates the hash", func() {
				pos.MakeMove(NewEnPassantMove(SQ_F5, SQ_G6))
				Expect(pos.hash).ToNot(Equal(prevHash))
			})
			It("updates a copy of frozenPos", func() {
				oldFP := pos.frozenPos
				pos.MakeMove(NewEnPassantMove(SQ_F5, SQ_G6))
				fp := pos.frozenPos
				Expect(fp).ToNot(Equal(oldFP))
				Expect(fp.Rule50).To(BeEquivalentTo(0))
				expCastleRights := [N_CASTLE_RIGHTS]bool{true, false, false, false}
				Expect(fp.CastleRights).To(Equal(expCastleRights))
				Expect(fp.EnPassantSq).To(Equal(NULL_SQ))
			})
			It("returns a black pawn", func() {
				capturedMove := pos.MakeMove(NewEnPassantMove(SQ_F5, SQ_G6))
				Expect(capturedMove).To(Equal(B_PAWN))
			})
		})
	})
	Describe("::MakeMove + ::UnmakeMove", func() {
		var pos *Position
		BeforeEach(func() {
			var posErr error
			pos, posErr = FromFEN("3k4/1b6/8/8/8/8/3P2B1/4K2R w - - 1 1")
			Expect(posErr).ToNot(HaveOccurred())
			expPieceBBs := [N_PIECES]Bitboard{
				0b11110111_11111101_11111111_11111111_11111111_11111111_10110111_01101111,
				BBWithSquares(SQ_D2),
				0,
				BBWithSquares(SQ_G2),
				BBWithSquares(SQ_H1),
				0,
				BBWithSquares(SQ_E1),
				0,
				0,
				BBWithSquares(SQ_B7),
				0,
				0,
				BBWithSquares(SQ_D8),
			}
			Expect(pos.pieceBitboards).To(Equal(expPieceBBs))
			expPieces := [N_SQUARES]Piece{}
			expPieces[SQ_E1] = W_KING
			expPieces[SQ_H1] = W_ROOK
			expPieces[SQ_D2] = W_PAWN
			expPieces[SQ_G2] = W_BISHOP
			expPieces[SQ_B7] = B_BISHOP
			expPieces[SQ_D8] = B_KING
			Expect(pos.pieces).To(Equal(expPieces))
			Expect(pos.colorBitboards[WHITE]).To(Equal(BBWithSquares(SQ_E1, SQ_H1, SQ_D2, SQ_G2)))
			Expect(pos.colorBitboards[BLACK]).To(Equal(BBWithSquares(SQ_B7, SQ_D8)))
		})
		When("the move is a pawn move", func() {
			It("restores the original position", func() {
				move := NewNormalMove(SQ_D2, SQ_D4)
				frozenPos := *pos.frozenPos
				capPiece := pos.MakeMove(move)
				pos.UnmakeMove(move, &frozenPos, capPiece)
				expPos, _ := FromFEN("3k4/1b6/8/8/8/8/3P2B1/4K2R w - - 1 1")
				Expect(pos).To(Equal(expPos))
			})
		})
		When("the move is castles", func() {
			It("restores the original position", func() {
				move := NewMove(SQ_E1, SQ_G1, NULL_SQ, EMPTY_PIECE_TYPE, true)
				frozenPos := *pos.frozenPos
				capPiece := pos.MakeMove(move)
				pos.UnmakeMove(move, &frozenPos, capPiece)
				expPos, _ := FromFEN("3k4/1b6/8/8/8/8/3P2B1/4K2R w - - 1 1")
				Expect(pos).To(Equal(expPos))
			})
		})
		When("the move is a capture", func() {
			It("restores the original position", func() {
				move := NewNormalMove(SQ_G2, SQ_B7)
				frozenPos := *pos.frozenPos
				capPiece := pos.MakeMove(move)
				pos.UnmakeMove(move, &frozenPos, capPiece)
				expPos, _ := FromFEN("3k4/1b6/8/8/8/8/3P2B1/4K2R w - - 1 1")
				Expect(pos).To(Equal(expPos))
			})
		})
		When("the move is en passant", func() {
			BeforeEach(func() {
				var posErr error
				pos, posErr = FromFEN("3k4/1b6/8/5Pp1/8/8/3P2B1/4K2R w K g6 1 1")
				Expect(posErr).ToNot(HaveOccurred())
				expPiecesBB := [N_PIECES]Bitboard{
					0b11110111_11111101_11111111_10011111_11111111_11111111_10110111_01101111,
					BBWithSquares(SQ_F5, SQ_D2),
					0,
					BBWithSquares(SQ_G2),
					BBWithSquares(SQ_H1),
					0,
					BBWithSquares(SQ_E1),
					BBWithSquares(SQ_G5),
					0,
					BBWithSquares(SQ_B7),
					0,
					0,
					BBWithSquares(SQ_D8),
				}
				Expect(pos.pieceBitboards).To(Equal(expPiecesBB))
			})
			It("restores the original position", func() {
				move := NewNormalMove(SQ_G2, SQ_B7)
				frozenPos := *pos.frozenPos
				capPiece := pos.MakeMove(move)
				pos.UnmakeMove(move, &frozenPos, capPiece)
				expPos, _ := FromFEN("3k4/1b6/8/5Pp1/8/8/3P2B1/4K2R w K g6 1 1")
				Expect(pos).To(Equal(expPos))
			})
		})
	})
})
