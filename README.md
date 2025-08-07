# YouTube Command-Line Tool

This is a powerful command-line tool for interacting with the YouTube Data API. It allows you to manage your YouTube channel, upload videos, and get detailed information about your content.

## Features

*   **List Channel Information:** Get basic information about a YouTube channel, including its ID, title, description, and view count.
*   **Upload Videos:** Upload videos to your YouTube channel with a variety of metadata options.
*   **List Videos:** Get a list of all the videos on a channel, with both basic and detailed views.
*   **Get Channel Info by Video ID:** Get the channel ID for a given video ID.
*   **Secure Authentication:** Uses OAuth 2.0 to securely authenticate with your YouTube account.
*   **Environment Variable Support:** Can be configured with environment variables for ease of use.

## Installation

1.  **Build the application:**

    ```bash
    go build
    ```

2.  **Set up your environment:**

    This tool requires a client secrets file from the Google Cloud Console. Once you have it, you can either provide the path to it with the `--secrets` flag, or you can set the `YOUTUBE_SECRETS` environment variable:

    ```bash
    export YOUTUBE_SECRETS=/path/to/your/client_secrets.json
    ```

## Usage

### `list`

Get basic information about a YouTube channel.

```bash
./mcp-youtube list --channelid <your-channel-id>
```

### `upload`

Upload a video to your YouTube channel.

```bash
./mcp-youtube upload --filename /path/to/your/video.mp4 --title "My Awesome Video" --description "This is a test video"
```

### `list-videos`

List the videos on a channel.

**Basic View:**

```bash
./mcp-youtube list-videos --channelid <your-channel-id>
```

**Detailed View:**

```bash
./mcp-youtube list-videos --channelid <your-channel-id> --detailed
```

**Pagination:**

The `list-videos` command supports pagination. You can use the `--limit` flag to control how many results are shown per page. You will be prompted to show more results if they are available.

### `channel-info`

Get the channel ID for a given video ID.

```bash
./mcp-youtube channel-info --videoid <video-id>
```
