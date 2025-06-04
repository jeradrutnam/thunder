#!/bin/bash
# ----------------------------------------------------------------------------
# Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.
# ----------------------------------------------------------------------------

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# --- Set Default OS and the architecture --- 
# Auto-detect GO OS
DEFAULT_OS=$(go env GOOS 2>/dev/null)
if [ -z "$DEFAULT_OS" ]; then
  UNAME_OS="$(uname -s)"
  case "$UNAME_OS" in
    Darwin) DEFAULT_OS="darwin" ;;
    Linux) DEFAULT_OS="linux" ;;
    MINGW*|MSYS*|CYGWIN*) DEFAULT_OS="windows" ;;
    *) echo "Unsupported OS: $UNAME_OS"; exit 1 ;;
  esac
fi
# Auto-detect GO ARCH
DEFAULT_ARCH=$(go env GOARCH 2>/dev/null)

echo $DEFAULT_ARCH
if [ -z "$DEFAULT_ARCH" ]; then
  UNAME_ARCH="$(uname -m)"
  case "$UNAME_ARCH" in
    x86_64|amd64) DEFAULT_ARCH="amd64" ;;
    arm64|aarch64) DEFAULT_ARCH="arm64" ;;
    *) echo "Unsupported architecture: $UNAME_ARCH"; exit 1 ;;
  esac
fi

GO_OS=${2:-$DEFAULT_OS}
GO_ARCH=${3:-$DEFAULT_ARCH}

SAMPLE_DIST_NODE_VERSION=node18
SAMPLE_DIST_OS=${2:-$DEFAULT_OS}
SAMPLE_DIST_ARCH=${3:-$DEFAULT_ARCH}

# Transform OS for node packaging executor
if [ "$SAMPLE_DIST_OS" = "darwin" ]; then
    SAMPLE_DIST_OS=macos
elif [ "$SAMPLE_DIST_OS" = "windows" ]; then
    SAMPLE_DIST_OS="win"
fi

if [ "$SAMPLE_DIST_ARCH" = "amd64" ]; then
    SAMPLE_DIST_ARCH=x64
fi

# --- Thunder Package Distribution details ---
GO_PACKAGE_OS=$GO_OS
GO_PACKAGE_ARCH=$GO_ARCH

# Normalize OS name for distribution packaging
if [ "$GO_OS" = "darwin" ]; then
    GO_PACKAGE_OS=macos
elif [ "$GO_OS" = "windows" ]; then
    GO_PACKAGE_OS="win"
fi

if [ "$GO_ARCH" = "amd64" ]; then
    GO_PACKAGE_ARCH=x64
fi

VERSION_FILE=version.txt
VERSION=$(cat "$VERSION_FILE")
BINARY_NAME=thunder
PRODUCT_FOLDER=${BINARY_NAME}-v${VERSION}-${GO_PACKAGE_OS}-${GO_PACKAGE_ARCH}

# --- Sample App Distribution details ---
SAMPLE_PACKAGE_OS=$SAMPLE_DIST_OS
SAMPLE_PACKAGE_ARCH=$SAMPLE_DIST_ARCH

SAMPLE_APP_SERVER_BINARY_NAME=server
SAMPLE_APP_VERSION=$(grep -o '"version": *"[^"]*"' samples/apps/oauth/package.json | sed 's/"version": *"\(.*\)"/\1/')
SAMPLE_APP_FOLDER="${BINARY_NAME}-sample-app-v${SAMPLE_APP_VERSION}-${SAMPLE_PACKAGE_OS}-${SAMPLE_PACKAGE_ARCH}"

# Server ports
BACKEND_PORT=8090

# Directories
TARGET_DIR=target
OUTPUT_DIR=$TARGET_DIR/out
DIST_DIR=$TARGET_DIR/dist
BUILD_DIR=$OUTPUT_DIR/.build
LOCAL_CERT_DIR=$OUTPUT_DIR/.cert
BACKEND_BASE_DIR=backend
BACKEND_DIR=$BACKEND_BASE_DIR/cmd/server
REPOSITORY_DIR=$BACKEND_BASE_DIR/cmd/server/repository
REPOSITORY_DB_DIR=$REPOSITORY_DIR/database
SERVER_SCRIPTS_DIR=$BACKEND_BASE_DIR/scripts
SERVER_DB_SCRIPTS_DIR=$BACKEND_BASE_DIR/dbscripts
SECURITY_DIR=repository/resources/security
SAMPLE_BASE_DIR=samples
SAMPLE_APP_DIR=$SAMPLE_BASE_DIR/apps/oauth
SAMPLE_APP_SERVER_DIR=$SAMPLE_APP_DIR/server

function clean_all() {
    echo "Cleaning all build artifacts..."
    rm -rf "$TARGET_DIR"

    echo "Removing certificates in the $BACKEND_DIR/$SECURITY_DIR"
    rm -rf "$BACKEND_DIR/$SECURITY_DIR"

    echo "Removing certificates in the $SAMPLE_APP_DIR"
    rm -f "$SAMPLE_APP_DIR/server.cert"
    rm -f "$SAMPLE_APP_DIR/server.key"
    rm -f "$SAMPLE_APP_SERVER_DIR/server.cert"
    rm -f "$SAMPLE_APP_SERVER_DIR/server.key"
}

function clean() {
    echo "Cleaning build artifacts..."
    rm -rf "$OUTPUT_DIR"

    echo "Removing certificates in the $BACKEND_DIR/$SECURITY_DIR"
    rm -rf "$BACKEND_DIR/$SECURITY_DIR"

    echo "Removing certificates in the $SAMPLE_APP_DIR"
    rm -f "$SAMPLE_APP_DIR/server.cert"
    rm -f "$SAMPLE_APP_DIR/server.key"
}

function build_backend() {
    echo "Building Go backend..."
    mkdir -p "$BUILD_DIR"

    # Set binary name with .exe extension for Windows
    local output_binary="$BINARY_NAME"
    if [ "$GO_OS" = "windows" ]; then
        output_binary="${BINARY_NAME}.exe"
    fi

    GO_OS=$GO_OS GO_ARCH=$GO_ARCH CGO_ENABLED=0 go build -C "$BACKEND_BASE_DIR" \
    -x -ldflags "-X \"main.version=$VERSION\" \
    -X \"main.buildDate=$$(date -u '+%Y-%m-%d %H:%M:%S UTC')\"" \
    -o "../$BUILD_DIR/$output_binary" ./cmd/server

    echo "Initializing databases..."
    initialize_databases
}

function initialize_databases() {
    echo "Initializing SQLite databases..."

    mkdir -p "$REPOSITORY_DB_DIR"

    db_files=("thunderdb.db" "runtimedb.db")
    script_paths=("thunderdb/sqlite.sql" "runtimedb/sqlite.sql")

    for ((i = 0; i < ${#db_files[@]}; i++)); do
        db_file="${db_files[$i]}"
        script_rel_path="${script_paths[$i]}"
        db_path="$REPOSITORY_DB_DIR/$db_file"
        script_path="$SERVER_DB_SCRIPTS_DIR/$script_rel_path"

        if [[ -f "$script_path" ]]; then
            if [[ -f "$db_path" ]]; then
                echo " - Removing existing $db_file"
                rm "$db_path"
            fi

            echo " - Creating $db_file using $script_path"
            sqlite3 "$db_path" < "$script_path"
        else
            echo " ! Skipping $db_file: SQL script not found at $script_path"
        fi
    done

    echo "SQLite database initialization complete."
}

function prepare_backend_for_packaging() {
    echo "Copying backend artifacts..."

    # Use appropriate binary name based on OS
    local binary_name="$BINARY_NAME"
    if [ "$GO_OS" = "windows" ]; then
        binary_name="${BINARY_NAME}.exe"
    fi

    cp "$BUILD_DIR/$binary_name" "$DIST_DIR/$PRODUCT_FOLDER/"
    cp -r "$REPOSITORY_DIR" "$DIST_DIR/$PRODUCT_FOLDER/"
    cp "$VERSION_FILE" "$DIST_DIR/$PRODUCT_FOLDER/"
    cp -r "$SERVER_SCRIPTS_DIR" "$DIST_DIR/$PRODUCT_FOLDER/"
    cp -r "$SERVER_DB_SCRIPTS_DIR" "$DIST_DIR/$PRODUCT_FOLDER/"
    mkdir -p "$DIST_DIR/$PRODUCT_FOLDER/$SECURITY_DIR"

    echo "=== Ensuring server certificates exist in the distribution ==="
    ensure_certificates "$DIST_DIR/$PRODUCT_FOLDER/$SECURITY_DIR"
}

function package_backend() {
    echo "Packaging backend artifacts..."

    mkdir -p "$DIST_DIR/$PRODUCT_FOLDER"

    prepare_backend_for_packaging

    # Copy the appropriate startup script based on the target OS
    if [ "$GO_OS" = "windows" ]; then
        echo "Including Windows start script (start.bat)..."
        cp -r "start.bat" "$DIST_DIR/$PRODUCT_FOLDER"
    else
        echo "Including Unix start script (start.sh)..."
        cp -r "start.sh" "$DIST_DIR/$PRODUCT_FOLDER"
    fi

    echo "Creating zip file..."
    (cd "$DIST_DIR" && zip -r "$PRODUCT_FOLDER.zip" "$PRODUCT_FOLDER")
    rm -rf "${DIST_DIR:?}/$PRODUCT_FOLDER" "$BUILD_DIR"
}

function build_sample_app() {
    echo "Building sample app..."

    # Ensure certificate exists for the sample app
    echo "=== Ensuring sample app certificates exist ==="
    ensure_certificates "$SAMPLE_APP_DIR"
    
    # Build the application
    cd "$SAMPLE_APP_DIR" || exit 1
    echo "Installing dependencies..."
    npm install
    
    echo "Building the app..."
    npm run build
    
    cd - || exit 1
    
    echo "Sample app built successfully."
}

function package_sample_app() {
    echo "Copying sample artifacts..."

    # Use appropriate binary name based on OS
    local binary_name="$SAMPLE_APP_SERVER_BINARY_NAME"
    if [ "$SAMPLE_DIST_OS" = "win" ]; then
        binary_name="${SAMPLE_APP_SERVER_BINARY_NAME}.exe"
    fi
    
    mkdir -p "$DIST_DIR/$SAMPLE_APP_FOLDER"
    
    # Copy the built app files
    cp -r "$SAMPLE_APP_SERVER_DIR/app" "$DIST_DIR/$SAMPLE_APP_FOLDER/"

    cd $SAMPLE_APP_SERVER_DIR

    mkdir -p "executables"

    npx pkg . -t $SAMPLE_DIST_NODE_VERSION-$SAMPLE_DIST_OS-$SAMPLE_DIST_ARCH -o executables/$SAMPLE_APP_SERVER_BINARY_NAME-$SAMPLE_DIST_OS-$SAMPLE_DIST_ARCH

    cd $SCRIPT_DIR

    # Copy the server binary
    cp "$SAMPLE_APP_SERVER_DIR/executables/$SAMPLE_APP_SERVER_BINARY_NAME-$SAMPLE_DIST_OS-$SAMPLE_DIST_ARCH" "$DIST_DIR/$SAMPLE_APP_FOLDER/$binary_name"

    # Ensure the certificates exist in the sample app directory
    echo "=== Ensuring certificates exist in the sample distribution ==="
    ensure_certificates "$DIST_DIR/$SAMPLE_APP_FOLDER"

    # Copy the appropriate startup script based on the target OS
    if [ "$SAMPLE_DIST_OS" = "win" ]; then
        echo "Including Windows start script (start.bat)..."
        cp -r "$SAMPLE_APP_SERVER_DIR/start.bat" "$DIST_DIR/$SAMPLE_APP_FOLDER"
    else
        echo "Including Unix start script (start.sh)..."
        cp -r "$SAMPLE_APP_SERVER_DIR/start.sh" "$DIST_DIR/$SAMPLE_APP_FOLDER"
    fi

    echo "Creating zip file..."
    (cd "$DIST_DIR" && zip -r "$SAMPLE_APP_FOLDER.zip" "$SAMPLE_APP_FOLDER")
    rm -rf "${DIST_DIR:?}/$SAMPLE_APP_FOLDER"
    
    echo "Sample app packaged successfully as $DIST_DIR/$SAMPLE_APP_FOLDER.zip"
}

function test_integration() {
    echo "Running integration tests..."
    go run -C ./tests/integration ./main.go
}

function ensure_certificates() {
    local cert_dir=$1
    local cert_name_prefix="server"
    local cert_file_name="${cert_name_prefix}.cert"
    local key_file_name="${cert_name_prefix}.key"

    # Generate certificate and key file if don't exists in the cert directory
    local local_cert_file="${LOCAL_CERT_DIR}/${cert_file_name}"
    local local_key_file="${LOCAL_CERT_DIR}/${key_file_name}"
    if [[ ! -f "$local_cert_file" || ! -f "$local_key_file" ]]; then
        mkdir -p "$LOCAL_CERT_DIR"
        echo "Generating SSL certificates in $LOCAL_CERT_DIR..."
        OPENSSL_ERR=$(
            openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
                -keyout "$local_key_file" \
                -out "$local_cert_file" \
                -subj "/O=WSO2/OU=Thunder/CN=localhost" \
                > /dev/null 2>&1
        )
        if [[ $? -ne 0 ]]; then
            echo "Error generating SSL certificates: $OPENSSL_ERR"
            exit 1
        fi
        echo "Certificates generated successfully in $LOCAL_CERT_DIR."
    else
        echo "Certificates already exist in $LOCAL_CERT_DIR."
    fi

    # Copy the generated certificates to the specified directory
    local cert_file="$cert_dir/${cert_file_name}"
    local key_file="$cert_dir/${key_file_name}"

    if [[ ! -f "$cert_file" || ! -f "$key_file" ]]; then
        mkdir -p "$cert_dir"
        echo "Copying certificates to $cert_dir..."
        cp "$local_cert_file" "$cert_file"
        cp "$local_key_file" "$key_file"
        echo "Certificates copied successfully to $cert_dir."
    else
        echo "Certificates already exist in $cert_dir."
    fi
}

function run() {
    echo "=== Ensuring server certificates exist ==="
    ensure_certificates "$BACKEND_DIR/$SECURITY_DIR"

    echo "=== Ensuring sample app certificates exist ==="
    ensure_certificates "$SAMPLE_APP_DIR"

    # Kill known ports
    function kill_port() {
        local port=$1
        lsof -ti tcp:$port | xargs kill -9 2>/dev/null || true
    }

    kill_port $BACKEND_PORT

    echo "=== Starting backend ==="
    BACKEND_PORT=$BACKEND_PORT go run -C "$BACKEND_DIR" . &
    BACKEND_PID=$!

    echo ""
    echo "⚡ Thunder Backend : https://localhost:$BACKEND_PORT"
    echo "Press Ctrl+C to stop."

    trap 'echo -e "\nStopping servers..."; kill $BACKEND_PID; exit' SIGINT

    wait $BACKEND_PID
}

case "$1" in
    clean)
        clean
        ;;
    clean_all)
        clean_all
        ;;
    build_backend)
        build_backend
        package_backend
        ;;
    build_samples)
        build_sample_app
        package_sample_app
        ;;
    build)
        build_backend
        package_backend
        build_sample_app
        package_sample_app
        ;;
    test)
        test_integration
        ;;
    run)
        run
        ;;
    *)
        echo "Usage: ./build.sh {clean|build|test|run} [OS] [ARCH]"
        echo ""
        echo "  clean         - Clean build artifacts"
        echo "  clean_all     - Clean all build artifacts including distributions"
        echo "  build         - Build the Thunder server only"
        echo "  build_backend - Build the Thunder backend server"
        echo "  build_samples - Build the sample applications"
        echo "  test          - Run integration tests"
        echo "  run           - Run the Thunder server for development"
        exit 1
        ;;
esac
