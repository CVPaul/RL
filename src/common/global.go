package common

const ( // global values for greedy_snake
	SEED          = 9527
	BODY_SIZE     = 20
	STATUS_X_SIZE = 20
	STATUS_Y_SIZE = 20
	MAGIC_NUM     = 2

	SLEEP_TIME         = 10
	TIME_FOR_SPEED_CTR = 100

	NOTHING_REWARD = -1
	FOOD_REWARD    = 100
	DEAD_LOSS      = -1000
	// actions
	ACTION_CNT = 4
	// parameters
	PARAM_POLICY_PATH = "data/policy."
	MINIMIZE_EPSILON  = 0.00
)

var ACTIONS = []string{"U", "D", "L", "R"}

const ( // window attributes
	WINDOW_WIDTH  = BODY_SIZE * STATUS_X_SIZE
	WINDOW_HEIGHT = BODY_SIZE * STATUS_Y_SIZE

	ADJ_MAT_SIZE = STATUS_X_SIZE * STATUS_Y_SIZE
	MAX_DIST     = ADJ_MAT_SIZE + 1000
)

const ( // metric
	NANOS_TO_MILLISECOND = 1000000
)

var GoStraight = map[string]string{"U": "U", "D": "D", "L": "L", "R": "R"}
var TurnLeft = map[string]string{"U": "L", "D": "R", "L": "D", "R": "U"}
var TurnRight = map[string]string{"U": "R", "D": "L", "L": "U", "R": "D"}

var GoStraight_r = map[string]string{"U": "U", "D": "D", "L": "L", "R": "R"}
var TurnLeft_r = map[string]string{"L": "U", "R": "D", "D": "L", "U": "R"}
var TurnRight_r = map[string]string{"R": "U", "L": "D", "U": "L", "D": "R"}

const ( // RL const define
	EAT_FOOD       = 550
	HIT_WALL       = -700
	DEFAULT_REWARD = -10
	DEFAULT_VALUE  = 0

	SAVE_EVERY_ITER = 1000000
	DEGREE_SPLIT    = 90
	EPSILON_FAC     = 1000000.0
	POWER_PARAM     = 1.0

	SHOW    = false
	DEBUG   = false
	VERBOSE = false
)
