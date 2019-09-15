# MKV-SEAT
MKV [S]ubtitle [E]xtract [A]nd [T]ranslate using Google Translation API.

The translation is optional, if you only want to extract the subtitle omit the flags `lang-from`, `lang-to` and `key`.

**You must have installed the mkvtoolnix** (https://mkvtoolnix.download/downloads.html) on your machine, otherwise it will not work.

## Usage
Checkout the releases page [https://github.com/herlon214/mkv-seat/releases](https://github.com/herlon214/mkv-seat/releases)
```
$ mkv-seat
          _                               _
_ __ ___ | | ____   __     ___  ___  __ _| |_
| '_ ` _ \| |/ /\ \ / /    / __|/ _ \/ _` | __|
| | | | | |   <  \ V /     \__ \  __/ (_| | |_
|_| |_| |_|_|\_\  \_/      |___/\___|\__,_|\__|


Error: requires at least 1 arg(s), only received 0
Usage:
  mkv-seat file.mkv [flags]

Flags:
  -h, --help                   help for mkv-seat
  -k, --key string             Google Cloud Translation Api Key, e.g: AIvaSyCiLjaWkykUoROHq2lqqbVoUA3ZTyv7xQI
  -f, --lang-from string       Original subtitle language (following the BCP 47), e.g: en
  -t, --lang-to string         Output subtitle language (following the BCP 47), e.g: pt-BR
  -o, --output-format string   Output format, e.g: srt (default "srt")
```

### Exctracting Subtitles Only
```
$ mkv-seat "[HorribleSubs] Kimetsu no Yaiba - 23 [1080p].mkv"
          _                               _
_ __ ___ | | ____   __     ___  ___  __ _| |_
| '_ ` _ \| |/ /\ \ / /    / __|/ _ \/ _` | __|
| | | | | |   <  \ V /     \__ \  __/ (_| | |_
|_| |_| |_|_|\_\  \_/      |___/\___|\__,_|\__|


INFO[0000] [MKV] Extracting subtitle for [HorribleSubs] Kimetsu no Yaiba - 23 [1080p]...
INFO[0000] [MKV] Subtitle extracted successfully
INFO[0000] Executed successfully!
```

### Extracting and Translating Subtitles

```
$ mkv-seat "[HorribleSubs] Kimetsu no Yaiba - 23 [1080p].mkv" -k YOUR_GOOGLE_APIS_KEY -f en -t pt-BR
          _                               _
_ __ ___ | | ____   __     ___  ___  __ _| |_
| '_ ` _ \| |/ /\ \ / /    / __|/ _ \/ _` | __|
| | | | | |   <  \ V /     \__ \  __/ (_| | |_
|_| |_| |_|_|\_\  \_/      |___/\___|\__,_|\__|


INFO[0000] [MKV] Extracting subtitle for [HorribleSubs] Kimetsu no Yaiba - 23 [1080p]...
INFO[0000] [MKV] Subtitle extracted successfully
INFO[0000] Translating subtitle from en to pt-BR
INFO[0000] [Translation] Collecting texts...
INFO[0000] [Translation] Collected 326 texts
INFO[0000] [Translation] Requesting translations 0/4...
INFO[0001] [Translation] Requesting translations 1/4...
INFO[0001] [Translation] Requesting translations 2/4...
INFO[0002] [Translation] Requesting translations 3/4...
INFO[0002] Executed successfully!
```