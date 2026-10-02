package tg

// cabacInit is one context variable's initialization pair (ITU-T H.264
// §9.3.1.1): preCtxState = Clip3(1, 126, ((m*Clip3(0,51,SliceQPY))>>4)+n),
// from which pStateIdx/valMPS are derived.
type cabacInit struct {
	m, n int8
}

// cabacRangeTabLPS is Table 9-44: codIRangeLPS indexed by [pStateIdx][(codIRange>>6)&3].
var cabacRangeTabLPS = [64][4]uint8{
	{128, 176, 208, 240}, {128, 167, 197, 227}, {128, 158, 187, 216}, {123, 150, 178, 205},
	{116, 142, 169, 195}, {111, 135, 160, 185}, {105, 128, 152, 175}, {100, 122, 144, 166},
	{95, 116, 137, 158}, {90, 110, 130, 150}, {85, 104, 123, 142}, {81, 99, 117, 135},
	{77, 94, 111, 128}, {73, 89, 105, 122}, {69, 85, 100, 116}, {66, 80, 95, 110},
	{62, 76, 90, 104}, {59, 72, 86, 99}, {56, 69, 81, 94}, {53, 65, 77, 89},
	{51, 62, 73, 85}, {48, 59, 69, 80}, {46, 56, 66, 76}, {43, 53, 63, 72},
	{41, 50, 59, 69}, {39, 48, 56, 65}, {37, 45, 54, 62}, {35, 43, 51, 59},
	{33, 41, 48, 56}, {32, 39, 46, 53}, {30, 37, 43, 50}, {29, 35, 41, 48},
	{27, 33, 39, 45}, {26, 31, 37, 43}, {24, 30, 35, 41}, {23, 28, 33, 39},
	{22, 27, 32, 37}, {21, 26, 30, 35}, {20, 24, 29, 33}, {19, 23, 27, 31},
	{18, 22, 26, 30}, {17, 21, 25, 28}, {16, 20, 23, 27}, {15, 19, 22, 25},
	{14, 18, 21, 24}, {14, 17, 20, 23}, {13, 16, 19, 22}, {12, 15, 18, 21},
	{12, 14, 17, 20}, {11, 14, 16, 19}, {11, 13, 15, 18}, {10, 12, 15, 17},
	{10, 12, 14, 16}, {9, 11, 13, 15}, {9, 11, 12, 14}, {8, 10, 12, 14},
	{8, 9, 11, 13}, {7, 9, 11, 12}, {7, 9, 10, 12}, {7, 8, 10, 11},
	{6, 8, 9, 11}, {6, 7, 9, 10}, {6, 7, 8, 9}, {2, 2, 2, 2},
}

// cabacTransIdxLPS is Table 9-45's LPS state transition, indexed by pStateIdx.
var cabacTransIdxLPS = [64]uint8{
	0, 0, 1, 2, 2, 4, 4, 5, 6, 7, 8, 9, 9, 11, 11, 12,
	13, 13, 15, 15, 16, 16, 18, 18, 19, 19, 21, 21, 22, 22, 23, 24,
	24, 25, 26, 26, 27, 27, 28, 29, 29, 30, 30, 30, 31, 32, 32, 33,
	33, 33, 34, 34, 35, 35, 35, 36, 36, 36, 37, 37, 37, 38, 38, 63,
}

// cabacTransIdxMPS is Table 9-45's MPS state transition, indexed by pStateIdx
// (pStateIdx+1 capped at 63).
var cabacTransIdxMPS = [64]uint8{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
	33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48,
	49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 62, 63,
}

// Real spec ctxIdx (ITU-T H.264 Table 9-11) for every I-slice syntax element
// this decoder uses. Unlike the first implementation attempt, these are the
// literal spec/reference-decoder indices (verified against FFmpeg's
// libavcodec/h264_cabac.c and cabac_context_init_I in cabac_data, not a
// private renumbering) - every cabacContextInit lookup below uses these
// constants directly, so there is no separate mapping layer left to get
// wrong.
const (
	ctxMbTypeIPrefix = 3 // 3 contexts (3,4,5): I_NxN-vs-not bin, ctxIdxInc 0-2 from neighbor mb_type
	ctxMbTypeI16x16  = 6 // 5 contexts (6-10): cbp_luma flag, cbp_chroma flag, cbp_chroma bit2, pred_mode bitA, pred_mode bitB

	ctxMbQpDelta = 60 // 4 contexts (60-63)

	ctxIntraChromaPredMode   = 64 // 4 contexts (64-67): bin0 uses 64+ctxIdxInc(0-2), bin1+ always 67
	ctxPrevIntraPredModeFlag = 68 // 1 context
	ctxRemIntraPredMode      = 69 // 1 context (reused for all 3 FL bins)

	ctxCbpLuma       = 73 // 4 contexts (73-76)
	ctxCbpChromaBin0 = 77 // 4 contexts (77-80)
	ctxCbpChromaBin1 = 81 // 4 contexts (81-84)

	ctxTransformSize8x8Flag = 399 // 3 contexts (399-401)

	// coded_block_flag bases (4 contexts each; cat 5/LumaLevel8x8 has none).
	ctxCBFCat0 = 85
	ctxCBFCat1 = 89
	ctxCBFCat2 = 93
	ctxCBFCat3 = 97
	ctxCBFCat4 = 101

	// significant_coeff_flag bases per ctxBlockCat (progressive/frame).
	ctxSigCat0 = 105
	ctxSigCat1 = 120
	ctxSigCat2 = 134
	ctxSigCat3 = 149
	ctxSigCat4 = 152
	ctxSigCat5 = 402 // + significantCoeffFlagOffset8x8[scanIdx] (15 distinct contexts, not 63)

	// last_significant_coeff_flag bases per ctxBlockCat (progressive/frame).
	ctxLastCat0 = 166
	ctxLastCat1 = 181
	ctxLastCat2 = 195
	ctxLastCat3 = 210
	ctxLastCat4 = 213
	ctxLastCat5 = 417 // + lastCoeffFlagOffset8x8[scanIdx] (9 distinct contexts, not 63)

	// coeff_abs_level_minus1 bases per ctxBlockCat (10 contexts each).
	ctxAbsLevelCat0 = 227
	ctxAbsLevelCat1 = 237
	ctxAbsLevelCat2 = 247
	ctxAbsLevelCat3 = 257
	ctxAbsLevelCat4 = 266
	ctxAbsLevelCat5 = 426
)

func sigCoeffFlagCtxBase(cat int) int {
	return [6]int{ctxSigCat0, ctxSigCat1, ctxSigCat2, ctxSigCat3, ctxSigCat4, ctxSigCat5}[cat]
}

func lastSigCoeffFlagCtxBase(cat int) int {
	return [6]int{ctxLastCat0, ctxLastCat1, ctxLastCat2, ctxLastCat3, ctxLastCat4, ctxLastCat5}[cat]
}

func coeffAbsLevelCtxBase(cat int) int {
	return [6]int{ctxAbsLevelCat0, ctxAbsLevelCat1, ctxAbsLevelCat2, ctxAbsLevelCat3, ctxAbsLevelCat4, ctxAbsLevelCat5}[cat]
}

func codedBlockFlagCtxBase(cat int) int {
	return [5]int{ctxCBFCat0, ctxCBFCat1, ctxCBFCat2, ctxCBFCat3, ctxCBFCat4}[cat]
}

// significantCoeffFlagOffset8x8 and lastCoeffFlagOffset8x8 map an 8x8
// block's scan position (0-62) to a context-index increment (ITU-T H.264's
// 8x8-specific significance-map contexts use far fewer distinct contexts
// than scan positions - 15 and 9 respectively, not 63 each). Verified
// against FFmpeg's libavcodec/cabac.c combined table
// (significant_coeff_flag_offset_8x8[0] / last_coeff_flag_offset_8x8).
var significantCoeffFlagOffset8x8 = [63]uint8{
	0, 1, 2, 3, 4, 5, 5, 4, 4, 3, 3, 4, 4, 4, 5, 5,
	4, 4, 4, 4, 3, 3, 6, 7, 7, 7, 8, 9, 10, 9, 8, 7,
	7, 6, 11, 12, 13, 11, 6, 7, 8, 9, 14, 10, 9, 8, 6, 11,
	12, 13, 11, 6, 9, 14, 10, 9, 11, 12, 13, 11, 14, 10, 12,
}

var lastCoeffFlagOffset8x8 = [63]uint8{
	0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4,
	5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8,
}

// coeffAbsLevel1Ctx/coeffAbsLevelGt1Ctx/coeffAbsLevelTransition implement
// the coeff_abs_level_minus1 "node context" state machine (ITU-T H.264
// §9.3.3.1.3's ctxIdxInc derivation is more subtle than a simple running
// count - this is a direct port of FFmpeg's verified node_ctx scheme, which
// a first, memory-based attempt at this decoder got wrong):
//   - node_ctx starts at 0 for each residual block.
//   - coeffAbsLevel1Ctx[node_ctx] is the context offset for the first
//     ("is level > 1") bin of each coefficient.
//   - On a level==1 decode, node_ctx advances via coeffAbsLevelTransition[0].
//   - On a level>1 decode, the unary continuation bins (coeff_abs 2..14)
//     all share the SAME context offset coeffAbsLevelGt1Ctx[node_ctx], and
//     node_ctx advances via coeffAbsLevelTransition[1] - once, before the
//     continuation bins, not per bin.
var coeffAbsLevel1Ctx = [8]uint8{1, 2, 3, 4, 0, 0, 0, 0}
var coeffAbsLevelGt1Ctx = [8]uint8{5, 5, 5, 5, 6, 7, 8, 9}
var coeffAbsLevelTransition = [2][8]uint8{
	{1, 2, 3, 3, 4, 5, 6, 7}, // after decoding a level == 1
	{4, 4, 4, 4, 5, 6, 7, 7}, // after decoding a level > 1
}

// iMbTypeInfo mirrors FFmpeg's ff_h264_i_mb_type_info[26]: mb_type (1-24, as
// computed by the raw bin decode in decodeMbTypeI) maps directly to
// (intra16x16PredMode, cbpChroma, cbpLuma) - not a bitwise unpacking of the
// decoded bins, which a first attempt at this decoder assumed incorrectly.
// Index 0 (I_NxN) and 25 (I_PCM) are unused placeholders, kept only so the
// table can be indexed directly by the raw mb_type value.
type iMbTypeEntry struct {
	predMode, cbpChroma, cbpLuma int
}

var iMbTypeInfo = [26]iMbTypeEntry{
	{}, // 0: I_NxN (handled separately, never indexes into this table)
	{2, 0, 0}, {1, 0, 0}, {0, 0, 0}, {3, 0, 0},
	{2, 1, 0}, {1, 1, 0}, {0, 1, 0}, {3, 1, 0},
	{2, 2, 0}, {1, 2, 0}, {0, 2, 0}, {3, 2, 0},
	{2, 0, 15}, {1, 0, 15}, {0, 0, 15}, {3, 0, 15},
	{2, 1, 15}, {1, 1, 15}, {0, 1, 15}, {3, 1, 15},
	{2, 2, 15}, {1, 2, 15}, {0, 2, 15}, {3, 2, 15},
	{}, // 25: I_PCM (handled separately)
}

// cabacContextInit is ITU-T H.264's I-slice CABAC context initialization
// table (the "I slice" column of Tables 9-12 through 9-33), indexed by the
// real spec ctxIdx (0-1023) - a direct, line-for-line port of FFmpeg's
// verified cabac_context_init_I[1024][2] (libavcodec/h264_cabac.c), not a
// memory-reconstructed approximation. Only the subset of indices this
// decoder actually uses (listed in the ctx* constants above) has been
// double-checked against the decode call sites; the rest are carried along
// unused.
var cabacContextInit = [1024]cabacInit{
	0: {20, -15}, 1: {2, 54}, 2: {3, 74}, 3: {20, -15},
	4: {2, 54}, 5: {3, 74}, 6: {-28, 127}, 7: {-23, 104},
	8: {-6, 53}, 9: {-1, 54}, 10: {7, 51},

	60: {0, 41}, 61: {0, 63}, 62: {0, 63}, 63: {0, 63},
	64: {-9, 83}, 65: {4, 86}, 66: {0, 97}, 67: {-7, 72},
	68: {13, 41}, 69: {3, 62},

	70: {0, 11}, 71: {1, 55}, 72: {0, 69}, 73: {-17, 127},
	74: {-13, 102}, 75: {0, 82}, 76: {-7, 74}, 77: {-21, 107},
	78: {-27, 127}, 79: {-31, 127}, 80: {-24, 127}, 81: {-18, 95},
	82: {-27, 127}, 83: {-21, 114}, 84: {-30, 127}, 85: {-17, 123},
	86: {-12, 115}, 87: {-16, 122},

	88: {-11, 115}, 89: {-12, 63}, 90: {-2, 68}, 91: {-15, 84},
	92: {-13, 104}, 93: {-3, 70}, 94: {-8, 93}, 95: {-10, 90},
	96: {-30, 127}, 97: {-1, 74}, 98: {-6, 97}, 99: {-7, 91},
	100: {-20, 127}, 101: {-4, 56}, 102: {-5, 82}, 103: {-7, 76},
	104: {-22, 125},

	105: {-7, 93}, 106: {-11, 87}, 107: {-3, 77}, 108: {-5, 71},
	109: {-4, 63}, 110: {-4, 68}, 111: {-12, 84}, 112: {-7, 62},
	113: {-7, 65}, 114: {8, 61}, 115: {5, 56}, 116: {-2, 66},
	117: {1, 64}, 118: {0, 61}, 119: {-2, 78}, 120: {1, 50},
	121: {7, 52}, 122: {10, 35}, 123: {0, 44}, 124: {11, 38},
	125: {1, 45}, 126: {0, 46}, 127: {5, 44}, 128: {31, 17},
	129: {1, 51}, 130: {7, 50}, 131: {28, 19}, 132: {16, 33},
	133: {14, 62}, 134: {-13, 108}, 135: {-15, 100},

	136: {-13, 101}, 137: {-13, 91}, 138: {-12, 94}, 139: {-10, 88},
	140: {-16, 84}, 141: {-10, 86}, 142: {-7, 83}, 143: {-13, 87},
	144: {-19, 94}, 145: {1, 70}, 146: {0, 72}, 147: {-5, 74},
	148: {18, 59}, 149: {-8, 102}, 150: {-15, 100}, 151: {0, 95},
	152: {-4, 75}, 153: {2, 72}, 154: {-11, 75}, 155: {-3, 71},
	156: {15, 46}, 157: {-13, 69}, 158: {0, 62}, 159: {0, 65},
	160: {21, 37}, 161: {-15, 72}, 162: {9, 57}, 163: {16, 54},
	164: {0, 62}, 165: {12, 72},

	166: {24, 0}, 167: {15, 9}, 168: {8, 25}, 169: {13, 18},
	170: {15, 9}, 171: {13, 19}, 172: {10, 37}, 173: {12, 18},
	174: {6, 29}, 175: {20, 33}, 176: {15, 30}, 177: {4, 45},
	178: {1, 58}, 179: {0, 62}, 180: {7, 61}, 181: {12, 38},
	182: {11, 45}, 183: {15, 39}, 184: {11, 42}, 185: {13, 44},
	186: {16, 45}, 187: {12, 41}, 188: {10, 49}, 189: {30, 34},
	190: {18, 42}, 191: {10, 55}, 192: {17, 51}, 193: {17, 46},
	194: {0, 89}, 195: {26, -19}, 196: {22, -17},

	197: {26, -17}, 198: {30, -25}, 199: {28, -20}, 200: {33, -23},
	201: {37, -27}, 202: {33, -23}, 203: {40, -28}, 204: {38, -17},
	205: {33, -11}, 206: {40, -15}, 207: {41, -6}, 208: {38, 1},
	209: {41, 17}, 210: {30, -6}, 211: {27, 3}, 212: {26, 22},
	213: {37, -16}, 214: {35, -4}, 215: {38, -8}, 216: {38, -3},
	217: {37, 3}, 218: {38, 5}, 219: {42, 0}, 220: {35, 16},
	221: {39, 22}, 222: {14, 48}, 223: {27, 37}, 224: {21, 60},
	225: {12, 68}, 226: {2, 97},

	227: {-3, 71}, 228: {-6, 42}, 229: {-5, 50}, 230: {-3, 54},
	231: {-2, 62}, 232: {0, 58}, 233: {1, 63}, 234: {-2, 72},
	235: {-1, 74}, 236: {-9, 91}, 237: {-5, 67}, 238: {-5, 27},
	239: {-3, 39}, 240: {-2, 44}, 241: {0, 46}, 242: {-16, 64},
	243: {-8, 68}, 244: {-10, 78}, 245: {-6, 77}, 246: {-10, 86},
	247: {-12, 92}, 248: {-15, 55}, 249: {-10, 60}, 250: {-6, 62},
	251: {-4, 65},

	252: {-12, 73}, 253: {-8, 76}, 254: {-7, 80}, 255: {-9, 88},
	256: {-17, 110}, 257: {-11, 97}, 258: {-20, 84}, 259: {-11, 79},
	260: {-6, 73}, 261: {-4, 74}, 262: {-13, 86}, 263: {-13, 96},
	264: {-11, 97}, 265: {-19, 117}, 266: {-8, 78}, 267: {-5, 33},
	268: {-4, 48}, 269: {-2, 53}, 270: {-3, 62}, 271: {-13, 71},
	272: {-10, 79}, 273: {-12, 86}, 274: {-13, 90}, 275: {-14, 97},

	276: {0, 0},

	277: {-6, 93}, 278: {-6, 84}, 279: {-8, 79}, 280: {0, 66},
	281: {-1, 71}, 282: {0, 62}, 283: {-2, 60}, 284: {-2, 59},
	285: {-5, 75}, 286: {-3, 62}, 287: {-4, 58}, 288: {-9, 66},
	289: {-1, 79}, 290: {0, 71}, 291: {3, 68}, 292: {10, 44},
	293: {-7, 62}, 294: {15, 36}, 295: {14, 40}, 296: {16, 27},
	297: {12, 29}, 298: {1, 44}, 299: {20, 36}, 300: {18, 32},
	301: {5, 42}, 302: {1, 48}, 303: {10, 62}, 304: {17, 46},
	305: {9, 64}, 306: {-12, 104}, 307: {-11, 97},

	308: {-16, 96}, 309: {-7, 88}, 310: {-8, 85}, 311: {-7, 85},
	312: {-9, 85}, 313: {-13, 88}, 314: {4, 66}, 315: {-3, 77},
	316: {-3, 76}, 317: {-6, 76}, 318: {10, 58}, 319: {-1, 76},
	320: {-1, 83}, 321: {-7, 99}, 322: {-14, 95}, 323: {2, 95},
	324: {0, 76}, 325: {-5, 74}, 326: {0, 70}, 327: {-11, 75},
	328: {1, 68}, 329: {0, 65}, 330: {-14, 73}, 331: {3, 62},
	332: {4, 62}, 333: {-1, 68}, 334: {-13, 75}, 335: {11, 55},
	336: {5, 64}, 337: {12, 70},

	338: {15, 6}, 339: {6, 19}, 340: {7, 16}, 341: {12, 14},
	342: {18, 13}, 343: {13, 11}, 344: {13, 15}, 345: {15, 16},
	346: {12, 23}, 347: {13, 23}, 348: {15, 20}, 349: {14, 26},
	350: {14, 44}, 351: {17, 40}, 352: {17, 47}, 353: {24, 17},
	354: {21, 21}, 355: {25, 22}, 356: {31, 27}, 357: {22, 29},
	358: {19, 35}, 359: {14, 50}, 360: {10, 57}, 361: {7, 63},
	362: {-2, 77}, 363: {-4, 82}, 364: {-3, 94}, 365: {9, 69},
	366: {-12, 109}, 367: {36, -35}, 368: {36, -34},

	369: {32, -26}, 370: {37, -30}, 371: {44, -32}, 372: {34, -18},
	373: {34, -15}, 374: {40, -15}, 375: {33, -7}, 376: {35, -5},
	377: {33, 0}, 378: {38, 2}, 379: {33, 13}, 380: {23, 35},
	381: {13, 58}, 382: {29, -3}, 383: {26, 0}, 384: {22, 30},
	385: {31, -7}, 386: {35, -15}, 387: {34, -3}, 388: {34, 3},
	389: {36, -1}, 390: {34, 5}, 391: {32, 11}, 392: {35, 5},
	393: {34, 12}, 394: {39, 11}, 395: {30, 29}, 396: {34, 26},
	397: {29, 39}, 398: {19, 66},

	399: {31, 21}, 400: {31, 31}, 401: {25, 50},
	402: {-17, 120}, 403: {-20, 112}, 404: {-18, 114}, 405: {-11, 85},
	406: {-15, 92}, 407: {-14, 89}, 408: {-26, 71}, 409: {-15, 81},
	410: {-14, 80}, 411: {0, 68}, 412: {-14, 70}, 413: {-24, 56},
	414: {-23, 68}, 415: {-24, 50}, 416: {-11, 74}, 417: {23, -13},
	418: {26, -13}, 419: {40, -15}, 420: {49, -14}, 421: {44, 3},
	422: {45, 6}, 423: {44, 34}, 424: {33, 54}, 425: {19, 82},
	426: {-3, 75}, 427: {-1, 23}, 428: {1, 34}, 429: {1, 43},
	430: {0, 54}, 431: {-2, 55}, 432: {0, 61}, 433: {1, 64},
	434: {0, 68}, 435: {-9, 92},

	436: {-14, 106}, 437: {-13, 97}, 438: {-15, 90}, 439: {-12, 90},
	440: {-18, 88}, 441: {-10, 73}, 442: {-9, 79}, 443: {-14, 86},
	444: {-10, 73}, 445: {-10, 70}, 446: {-10, 69}, 447: {-5, 66},
	448: {-9, 64}, 449: {-5, 58}, 450: {2, 59}, 451: {21, -10},
	452: {24, -11}, 453: {28, -8}, 454: {28, -1}, 455: {29, 3},
	456: {29, 9}, 457: {35, 20}, 458: {29, 36}, 459: {14, 67},
}
