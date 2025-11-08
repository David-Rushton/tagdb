[CmdletBinding()]
param (
    [parameter(Mandatory=$true)]
    [string]$Version
)

$tag = "tagdb-ws:$Version"

try {
    # Build the binaries.
    $env:GOOS = "linux"
    go build -v ./cmd/tagdb_ws
    go build -v ./cmd/tagdb_cli

    # Build the Docker image.
    docker build --tag $tag .

    # Cleanup any existing containers.
    try {
        docker rm --force tagdb-ws
    } catch {
        if ($_ -notmatch "No such container") {
            throw $_
        }
    }

    # Run the Docker container.
    docker run `
        --name tagdb-ws `
        --detach `
        --mount type=bind,src=C:\Users\davidr\.q,dst=/data `
        --publish 31979:8080 `
        --restart unless-stopped `
        --env TAGDB_PORT=8080 `
        --env TAGDB_WEB_ROOT=/web `
        --env TAGDB_STORAGE_ROOT=/data `
        --env TAGDB_STORAGE_WAL_ROLL_AFTER_BYTES=10485760 `
        --env TAGDB_STORAGE_BACKGROUND_TASK_INTERVAL_MS=5000 `
        $tag
}
catch {
    throw $_
}
finally {
    # Make sure to clean up the environment variable.
    Remove-Item env:GOOS

    # Clean up any built binaries.
    Remove-Item tagdb_ws
    Remove-Item tagdb_cli
}
