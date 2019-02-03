# Discord Music Bot

## **Modules:**

### **Bot**

Bot is the main module. It receives commands from a Discord channel and executes them.

### **Player**

Player handles audio playback

### **Downloader**

Downloader handles downloading of audio using a queue

### **Audio repository**

Audio repository handles storing audio data temporarily in the form of bytes or a file

### **Videoaudio**

Videoaudio handles encoding and transcoding

### **Youtube**

Youtube basically wraps a Youtube library only for the needed functionality

### **Spotify**

Spotify accesses the Spotify API which is pretty much unused for the time being

### **Soundboard**

Soundboard binary (not in git) plays a sound file through speakers or writes the data to standard output. It's used for reading sound files that'll be streamed to Discord. Streaming not implemented for now.

***
#### Environment variables:

- BOT_TOKEN - Discord Bot token
- CLIENT_ID - Discord API Client ID
- SPOTIFY_ID - Spotify API Client ID
- SPOTIFY_SECRET - Spotify API Client Secret

Environment variables are put into the environment from <i>.env</i> file.

***

### Configuration file - config.yaml

debug: boolean
