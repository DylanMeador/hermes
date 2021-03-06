package sounds

type Sound string

const (
	BLAME                    = Sound("pkg/sounds/blame.dca")
	HYPE1                    = Sound("pkg/sounds/god-that-looks-so-clean.dca")
	HYPE2                    = Sound("pkg/sounds/youre-the-best-player-in-the-game.dca")
	HORNHORNHORN             = Sound("pkg/sounds/mlg.dca")
	SHACO_ATTACK1            = Sound("pkg/sounds/shaco/attack1.dca")
	SHACO_ATTACK2            = Sound("pkg/sounds/shaco/attack2.dca")
	SHACO_JOKES_ON_YOU       = Sound("pkg/sounds/shaco/attack3.dca")
	SHACO_ATTACK4            = Sound("pkg/sounds/shaco/attack4.dca")
	SHACO_ATTACK5            = Sound("pkg/sounds/shaco/attack5.dca")
	SHACO_ATTACK6            = Sound("pkg/sounds/shaco/attack6.dca")
	SHACO_ATTACK7            = Sound("pkg/sounds/shaco/attack7.dca")
	SHACO_MAKE_YOU_DISAPPEAR = Sound("pkg/sounds/shaco/joke.dca")
	SHACO_LAUGH1             = Sound("pkg/sounds/shaco/laugh1.dca")
	SHACO_LAUGH2             = Sound("pkg/sounds/shaco/laugh2.dca")
	SHACO_LAUGH3             = Sound("pkg/sounds/shaco/laugh3.dca")
	SHACO_MAGIC_TRICK        = Sound("pkg/sounds/shaco/select.dca")
	SHACO_TAUNT              = Sound("pkg/sounds/shaco/taunt.dca")
	SAD                      = Sound("pkg/sounds/imogen-heap-hide-and-seek-mm-what-you-say.dca")
	RECORDING_PROMPT         = Sound("pkg/sounds/recording.dca")
	COUNTDOWN                = Sound("pkg/sounds/genji54321.dca")
)

var (
	ALL_SHACO = []Sound{
		SHACO_ATTACK1, SHACO_ATTACK2, SHACO_JOKES_ON_YOU, SHACO_ATTACK4, SHACO_ATTACK5, SHACO_ATTACK6, SHACO_ATTACK7,
		SHACO_MAKE_YOU_DISAPPEAR, SHACO_LAUGH1, SHACO_LAUGH2, SHACO_LAUGH3, SHACO_MAGIC_TRICK, SHACO_TAUNT,
	}
)
