package commands

const (
	COMMAND_PREFIX = "!"
	PAUSED_PREFIX  = "[Paused]"
)

// JOIN_DEFAULT_CHANNEL joins the default voice channel
const JOIN_DEFAULT_CHANNEL = "start"

// LEAVE_CHANNEL leaves the voice channel
const LEAVE_CHANNEL = "stop"

// JOIN_CHANNEL joins a specific voice channel
const JOIN_CHANNEL = "join "

// REPEAT_TEXT repeats the text
const REPEAT_TEXT = "repeat "

// FIND_PLAYLIST searches for spotify Playlists
const FIND_PLAYLIST = "playlist "

// PLAY_TRACK queues a track
const PLAY_TRACK = "play "

// UNPAUSE unpauses playback
const UNPAUSE = "play"

// PAUSE pauses playback
const PAUSE = "pause"

// SOUND plays a sound clip from the soundboard
const SOUND = "sound "
